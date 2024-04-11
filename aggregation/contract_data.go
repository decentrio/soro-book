package aggregation

import (
	"fmt"
	"time"

	"github.com/decentrio/soro-book/database/models"
)

func (as *Aggregation) contractDataEntryProcessing() {
	for {
		// Block until state have sync successful
		if as.isReSync {
			continue
		}

		select {
		// Receive a new tx
		case e := <-as.contractDataEntrysQueue:
			as.handleReceiveNewContractDataEntry(e)
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		}
		time.Sleep(time.Millisecond)
	}
}

func (as *Aggregation) handleReceiveNewContractDataEntry(e models.Contract) {
	_, err := as.db.CreateContractEntry(&e)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error create contract data entry ledger %d tx %s: %s", e.Ledger, e.TxHash, err.Error()))
	}
}

func (tw TransactionWrapper) GetModelsContractDataEntry() []models.Contract {
	v3 := tw.Tx.UnsafeMeta.V3
	if v3 == nil {
		return nil
	}

	var entries []models.Contract
	for _, op := range v3.Operations {
		for _, change := range op.Changes {
			entry, entryType, found := ContractDataEntry(change)
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
					TxHash:     tw.GetTransactionHash(),
					AccountId:  accountId,
					Ledger:     tw.GetLedgerSequence(),
					EntryType:  entryType,
					KeyXdr:     keyBz,
					ValueXdr:   valBz,
					Durability: int32(entry.Durability),
					IsNewest:   true,
				}
				entries = append(entries, entry)
			}
		}
	}

	return entries
}
