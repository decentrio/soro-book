package aggregation

import (
	"fmt"
	"io"
	"time"

	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/ingest"
	backends "github.com/stellar/go/ingest/ledgerbackend"
	"github.com/stellar/go/xdr"
)

type LedgerWrapper struct {
	ledger models.Ledger
	txs    []TransactionWrapper
}

func (as *Aggregation) getNewLedger() {
	from := as.startLedgerSeq
	to := as.startLedgerSeq + as.prepareStep
	ledgerRange := backends.BoundedRange(from, to)
	err := as.backend.PrepareRange(as.ctx, ledgerRange)
	if err != nil {
		//"is greater than max available in history archives"
		fmt.Println("Error Prepare")
		if as.prepareStep > 1 {
			as.prepareStep = as.prepareStep / 2
		} else {
			as.prepareStep = 1
		}
		time.Sleep(time.Millisecond)
		return
	}
	for seq := from; seq < to; seq++ {
		// get ledger
		ledgerCloseMeta, err := as.backend.GetLedger(as.ctx, seq)
		if err != nil {
			as.Logger.Error(fmt.Sprintf("Error GetLedger %s", err.Error()))
			continue
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
				as.Logger.Error(fmt.Sprintf("Error txReader %s", err.Error()))
			}

			txWrapper := NewTransactionWrapper(tx, seq)
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
	}
	as.startLedgerSeq = to
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
			as.Logger.Info("getting new ledgers")
			as.handleReceiveNewLedger(ledger)

			as.state = TX
			as.Logger.Info(fmt.Sprintf("state transition state: %s", as.state))

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
	_, err := as.db.CreateLedger(&lw.ledger)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error create ledger %d: %s", lw.ledger.Seq, err.Error()))
	}

	// Create Tx and Soroban events
	for _, tw := range lw.txs {
		as.txQueue <- tw
	}
}

func getLedgerFromCloseMeta(ledgerCloseMeta xdr.LedgerCloseMeta) models.Ledger {
	return models.Ledger{
		Hash:     ledgerCloseMeta.LedgerHash().HexString(),
		PrevHash: ledgerCloseMeta.PreviousLedgerHash().HexString(),
		Seq:      ledgerCloseMeta.LedgerSequence(),
	}
}
