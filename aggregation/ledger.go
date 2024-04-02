package aggregation

import (
	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/xdr"
)

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
