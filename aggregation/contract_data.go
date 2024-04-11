package aggregation

import (
	"fmt"
	"time"

	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/strkey"
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
			// continue with "state" because we don't want to store this entry
			if entryType == "state" {
				continue
			}
			if found {
				keyBz, _ := entry.Key.MarshalBinary()
				valBz, _ := entry.Val.MarshalBinary()
				var contractId string
				var err error
				if entry.Contract.ContractId != nil {
					contractId, err = strkey.Encode(strkey.VersionByteContract, entry.Contract.ContractId[:])
					if err != nil {
						continue
					}
				}

				var accountId string
				if entry.Contract.AccountId != nil {
					accountId, err = entry.Contract.AccountId.GetAddress()
					if err != nil {
						continue
					}
				}

				entry := models.Contract{
					ContractId: contractId,
					AccountId:  accountId,
					TxHash:     tw.GetTransactionHash(),
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
