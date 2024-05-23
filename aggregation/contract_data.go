package aggregation

import (
	"fmt"
	"time"

	"github.com/decentrio/soro-book/database/models"
	"github.com/google/uuid"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

func (as *Aggregation) contractDataEntryProcessing() {
	for {
		if as.state != CONTRACT {
			continue
		}

		select {
		// Receive a new tx
		case e := <-as.contractDataEntrysQueue:
			as.Logger.Info("getting new contract data entry")
			as.handleReceiveNewContractDataEntry(e)
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}
}

func (as *Aggregation) handleReceiveNewContractDataEntry(e models.ContractsData) {
	_, err := as.db.CreateContractEntry(&e)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error create contract data entry ledger %d tx %s: %s", e.Ledger, e.TxHash, err.Error()))
	}
}

func (tw TransactionWrapper) GetModelsContractDataEntry() []models.ContractsData {
	v3 := tw.Tx.UnsafeMeta.V3
	if v3 == nil {
		return nil
	}

	var entries []models.ContractsData
	for _, op := range v3.Operations {
		for _, change := range op.Changes {
			entry, entryType, found := ContractDataEntry(change)
			// continue with "state" because we don't want to store this entry
			if entryType == "state" {
				continue
			}
			if found {
				keyBz, _ := entry.Key.MarshalBinary()

				var valBz []byte
				if entryType != "removed" {
					valBz, _ = entry.Val.MarshalBinary()
				}
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

				entry := models.ContractsData{
					Id:         uuid.New().String(),
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

func ContractDataEntry(c xdr.LedgerEntryChange) (xdr.ContractDataEntry, string, bool) {
	var result xdr.ContractDataEntry

	switch c.Type {
	case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
		created := *c.Created
		if created.Data.ContractData != nil {
			result = *created.Data.ContractData
			return result, "created", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
		updated := *c.Updated
		if updated.Data.ContractData != nil {
			result = *updated.Data.ContractData
			return result, "updated", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
		ledgerKey := c.Removed
		if ledgerKey.ContractData != nil {
			result.Contract = ledgerKey.ContractData.Contract
			result.Key = ledgerKey.ContractData.Key
			result.Durability = ledgerKey.ContractData.Durability
			return result, "removed", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryState:
		state := *c.State
		if state.Data.ContractData != nil {
			result = *state.Data.ContractData
			return result, "state", true
		}

	}
	return result, "", false
}
