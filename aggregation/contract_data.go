package aggregation

import (
	"github.com/decentrio/soro-book/database/models"
)

func (tw TransactionWrapper) GetModelsContractDataEntry() []models.Contract {
	v3 := tw.Tx.UnsafeMeta.V3
	if v3 == nil {
		return nil
	}

	var entries []models.Contract
	for _, op := range v3.Operations {
		for _, change := range op.Changes {
			entry, found := ContractDataEntry(change)
			if found {
				keyBz, _ := entry.Key.MarshalBinary()
				valBz, _ := entry.Val.MarshalBinary()
				var contractId string
				if entry.Contract.ContractId != nil {
					contractId = (*entry.Contract.ContractId).HexString()
				}

				var accountId string
				if entry.Contract.AccountId != nil {
					accountId = (*entry.Contract.AccountId).Address()
				}

				entry := models.Contract{
					ContractId: contractId,
					AccountId:  accountId,
					Ledger:     tw.GetLedgerSequence(),
					KeyXdr:     keyBz,
					ValueXdr:   valBz,
					Durability: int32(entry.Durability),
				}
				entries = append(entries, entry)
			}
		}
	}

	return entries
}
