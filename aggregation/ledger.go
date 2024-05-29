package aggregation

import (
	"fmt"
	"io"
	"time"

	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
)

type LedgerWrapper struct {
	ledger models.Ledger
	txs    []TransactionWrapper
}

func (as *Aggregation) getNewLedger() {
	// get ledger
	seq := as.startLedgerSeq
	ledgerCloseMeta, err := as.backend.GetLedger(as.ctx, seq)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("error get ledger %s", err.Error()))
		return
	}

	ledger := getLedgerFromCloseMeta(ledgerCloseMeta)

	var txWrappers []TransactionWrapper
	var transactions = uint32(0)
	var operations = uint32(0)
	// get tx
	txReader, err := ingest.NewLedgerTransactionReader(
		as.ctx, as.backend, as.Cfg.NetworkPassphrase, seq,
	)
	panicIf(err)
	defer txReader.Close()

	// Read each transaction within the ledger, extract its operations, and
	// accumulate the statistics we're interested in.
	for {
		tx, err := txReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			as.Logger.Error(fmt.Sprintf("error txReader %s", err.Error()))
		}

		txWrapper := NewTransactionWrapper(tx, seq, ledger.LedgerTime)
		txWrappers = append(txWrappers, txWrapper)

		operations += uint32(len(tx.Envelope.Operations()))
		transactions++
	}

	ledger.Transactions = transactions
	ledger.Operations = operations

	lw := LedgerWrapper{
		ledger: ledger,
		txs:    txWrappers,
	}
	as.ledgerQueue <- lw
	as.startLedgerSeq++
}

// aggregation process
func (as *Aggregation) ledgerProcessing() {
	for {
		if as.state != LEDGER {
			continue
		}

		select {
		// Receive a new tx
		case ledger := <-as.ledgerQueue:
			as.handleReceiveNewLedger(ledger)

			as.state = TX

			as.CurrLedgerSeq = ledger.ledger.Seq
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}
}

// handleReceiveTx
func (as *Aggregation) handleReceiveNewLedger(lw LedgerWrapper) {
	// Create Ledger
	// _, err := as.db.CreateLedger(&lw.ledger)
	// if err != nil {
	// 	as.Logger.Error(fmt.Sprintf("Error create ledger %d: %s", lw.ledger.Seq, err.Error()))
	// }

	// Create Tx and Soroban events
	for _, tw := range lw.txs {
		_ = tw
		// as.txQueue <- tw
	}
}

func getLedgerFromCloseMeta(ledgerCloseMeta xdr.LedgerCloseMeta) models.Ledger {
	var ledgerHeader xdr.LedgerHeaderHistoryEntry
	switch ledgerCloseMeta.V {
	case 0:
		ledgerHeader = ledgerCloseMeta.MustV0().LedgerHeader
	case 1:
		ledgerHeader = ledgerCloseMeta.MustV1().LedgerHeader
	default:
		panic(fmt.Sprintf("Unsupported LedgerCloseMeta.V: %d", ledgerCloseMeta.V))
	}

	timeStamp := uint64(ledgerHeader.Header.ScpValue.CloseTime)

	return models.Ledger{
		Hash:       ledgerCloseMeta.LedgerHash().HexString(),
		PrevHash:   ledgerCloseMeta.PreviousLedgerHash().HexString(),
		Seq:        ledgerCloseMeta.LedgerSequence(),
		LedgerTime: timeStamp,
	}
}
