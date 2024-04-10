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

func (as *Aggregation) getNewLedger() {
	from := as.sequence
	to := as.sequence + prepareStep
	ledgerRange := backends.BoundedRange(from, to)
	err := as.backend.PrepareRange(as.ctx, ledgerRange)
	if err != nil {
		//"is greater than max available in history archives"
		time.Sleep(time.Second)
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
			as.ctx, as.backend, Config.NetworkPassphrase, seq,
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

		go func(lwi LedgerWrapper) {
			as.ledgerQueue <- lwi
		}(lw)
	}
	as.sequence = to
}

// aggregation process
func (as *Aggregation) ledgerProcessing() {
	for {
		// Block until state have sync successful
		if as.isReSync {
			continue
		}

		select {
		// Receive a new tx
		case ledger := <-as.ledgerQueue:
			as.handleReceiveNewLedger(ledger)
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		}
		time.Sleep(time.Millisecond)
	}
}

// handleReceiveTx
func (as *Aggregation) handleReceiveNewLedger(lw LedgerWrapper) {
	// Create Ledger
	_, err := as.db.CreateLedger(&lw.ledger)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error create ledger %d: %s", lw.ledger.Sequence, err.Error()))
	}

	// Create Tx and Soroban events
	for _, tw := range lw.txs {
		go func(twi TransactionWrapper) {
			as.txQueue <- twi
		}(tw)
	}
}

func getLedgerFromCloseMeta(ledgerCloseMeta xdr.LedgerCloseMeta) models.Ledger {
	return models.Ledger{
		Hash:     ledgerCloseMeta.LedgerHash().HexString(),
		PrevHash: ledgerCloseMeta.PreviousLedgerHash().HexString(),
		Sequence: ledgerCloseMeta.LedgerSequence(),
	}
}

func ContractDataEntry(c xdr.LedgerEntryChange) (xdr.ContractDataEntry, bool) {
	var result xdr.ContractDataEntry

	switch c.Type {
	case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
		created := *c.Created
		if created.Data.ContractData != nil {
			result = *created.Data.ContractData
			return result, true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
		updated := *c.Updated
		if updated.Data.ContractData != nil {
			result = *updated.Data.ContractData
			return result, true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
		return result, false
	case xdr.LedgerEntryChangeTypeLedgerEntryState:
		state := *c.State
		if state.Data.ContractData != nil {
			result = *state.Data.ContractData
			return result, true
		}

	}
	return result, false
}
