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
	// prepare range
	from, to := as.prepare()
	// get ledger
	if !as.isSync {
		for seq := from; seq < to; seq++ {
			ledgerCloseMeta, err := as.backend.GetLedger(as.ctx, seq)
			if err != nil {
				as.Logger.Error(fmt.Sprintf("error get ledger %s", err.Error()))
				return
			}

			go func(l xdr.LedgerCloseMeta) {
				as.ledgerQueue <- l
			}(ledgerCloseMeta)
		}
	} else {
		seq := as.StartLedgerSeq
		ledgerCloseMeta, err := as.backend.GetLedger(as.ctx, seq)
		if err != nil {
			as.Logger.Error(fmt.Sprintf("error get ledger %s", err.Error()))
			return
		}

		go func(l xdr.LedgerCloseMeta) {
			as.ledgerQueue <- l
		}(ledgerCloseMeta)
		as.StartLedgerSeq++
	}
}

// aggregation process
func (as *Aggregation) ledgerProcessing() {
	for {
		select {
		// Receive a new tx
		case ledger := <-as.ledgerQueue:
			as.handleReceiveNewLedger(ledger)
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}
}

// handleReceiveTx
func (as *Aggregation) handleReceiveNewLedger(l xdr.LedgerCloseMeta) {
	ledger := getLedgerFromCloseMeta(l)

	var txWrappers []TransactionWrapper
	var transactions = uint32(0)
	var operations = uint32(0)
	// get tx
	txReader, err := ingest.NewLedgerTransactionReaderFromLedgerCloseMeta(as.Cfg.NetworkPassphrase, l)
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

		txWrapper := NewTransactionWrapper(tx, ledger.Seq, ledger.LedgerTime)
		txWrappers = append(txWrappers, txWrapper)

		operations += uint32(len(tx.Envelope.Operations()))
		transactions++
	}

	ledger.Transactions = transactions
	ledger.Operations = operations

	// Create Ledger
	_, err = as.db.CreateLedger(&ledger)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error create ledger %d: %s", ledger.Seq, err.Error()))
	}

	// Create Tx and Soroban events
	for _, tw := range txWrappers {
		go func(twi TransactionWrapper) {
			as.txQueue <- twi
		}(tw)
	}
}

func (as *Aggregation) prepare() (uint32, uint32) {
	if !as.isSync {
		from := as.StartLedgerSeq
		to := from + DefaultPrepareStep

		var ledgerRange backends.Range
		if to > as.CurrLedgerSeq {
			ledgerRange = backends.UnboundedRange(from)
		} else {
			ledgerRange = backends.BoundedRange(from, to)
		}

		fmt.Println(ledgerRange)
		err := as.backend.PrepareRange(as.ctx, ledgerRange)
		if err != nil {
			as.Logger.Errorf("error prepare %s", err.Error())
			return 0, 0 // if prepare error, we should skip here
		} else {
			if to > as.CurrLedgerSeq {
				as.isSync = true
			}
		}
		as.StartLedgerSeq += DefaultPrepareStep
		return from, to
	}

	return 0, 0
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
