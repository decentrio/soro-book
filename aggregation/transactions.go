package aggregation

import (
	"fmt"

	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
)

const (
	SUCCESS = "success"
	FAILED  = "failed"
)

func (as *Aggregation) transactionProcessing() {
	for {
		// Block until state have sync successful
		if as.isReSync {
			continue
		}

		select {
		// Receive a new tx
		case tx := <-as.txQueue:
			as.handleReceiveNewTransaction(tx)
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		}
	}
}

func (as *Aggregation) handleReceiveNewTransaction(tw TransactionWrapper) {
	tx := tw.GetModelsTransaction()
	_, err := as.db.CreateTransaction(tx)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error create ledger %d tx %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), err.Error()))
	}

	// Contract entry
	// entries := tw.GetModelsContractDataEntry()
	// for _, entry := range entries {
	// 	_, err := as.db.CreateContractEntry(&entry)
	// 	if err != nil {
	// 		as.Logger.Error(fmt.Sprintf("Error create contract data entry ledger %d tx %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), err.Error()))
	// 		continue
	// 	}
	// }

	// Check if tx metadata is v3
	// txMetaV3, ok := tw.Tx.UnsafeMeta.GetV3()
	// if !ok {
	// 	continue
	// }

	// if txMetaV3.SorobanMeta == nil {
	// 	continue
	// }

	// // Create Event
	// for _, op := range tw.Ops {
	// 	events := op.GetContractEvents()
	// 	for _, event := range events {
	// 		_, err := as.db.CreateWasmContractEvent(&event)
	// 		if err != nil {
	// 			as.Logger.Error(fmt.Sprintf("Error create event ledger %d tx %s event %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), event.ContractId, err.Error()))
	// 			continue
	// 		}
	// 	}
	// }
}

type TransactionWrapper struct {
	LedgerSequence uint32
	Tx             ingest.LedgerTransaction
	Ops            []transactionOperationWrapper
}

func NewTransactionWrapper(tx ingest.LedgerTransaction, seq uint32) TransactionWrapper {
	var ops []transactionOperationWrapper
	for opi, op := range tx.Envelope.Operations() {
		operation := transactionOperationWrapper{
			index:          uint32(opi),
			txIndex:        tx.Index,
			operation:      op,
			ledgerSequence: seq,
		}

		ops = append(ops, operation)
	}

	return TransactionWrapper{
		LedgerSequence: seq,
		Tx:             tx,
		Ops:            ops,
	}
}

func (tw TransactionWrapper) GetTransactionHash() string {
	return tw.Tx.Result.TransactionHash.HexString()
}

func (tw TransactionWrapper) GetStatus() string {
	if tw.Tx.Result.Successful() {
		return SUCCESS
	}

	return FAILED
}

func (tw TransactionWrapper) GetLedgerSequence() uint32 {
	return tw.LedgerSequence
}

func (tw TransactionWrapper) GetApplicationOrder() uint32 {
	return tw.Tx.Index
}

func (tw TransactionWrapper) GetEnvelopeXdr() []byte {
	bz, _ := tw.Tx.Envelope.MarshalBinary()
	return bz
}

func (tw TransactionWrapper) GetResultXdr() []byte {
	bz, _ := tw.Tx.Result.MarshalBinary()
	return bz
}

func (tw TransactionWrapper) GetResultMetaXdr() []byte {
	txResultMeta := xdr.TransactionResultMeta{
		Result:            tw.Tx.Result,
		FeeProcessing:     tw.Tx.FeeChanges,
		TxApplyProcessing: tw.Tx.UnsafeMeta,
	}

	bz, _ := txResultMeta.MarshalBinary()

	return bz
}

func (tw TransactionWrapper) GetModelsTransaction() *models.Transaction {
	return &models.Transaction{
		Hash:             tw.GetTransactionHash(),
		Status:           tw.GetStatus(),
		Ledger:           tw.GetLedgerSequence(),
		ApplicationOrder: tw.GetApplicationOrder(),
		EnvelopeXdr:      tw.GetEnvelopeXdr(),   // xdr.TransactionEnvelope
		ResultXdr:        tw.GetResultXdr(),     // xdr.TransactionResultPair
		ResultMetaXdr:    tw.GetResultMetaXdr(), //xdr.TransactionResultMeta
		SourceAddress:    tw.Tx.Envelope.SourceAccount().ToAccountId().Address(),
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
