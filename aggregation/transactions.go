package aggregation

import (
	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/xdr"
)

const (
	SUCCESS = "success"
	FAILED  = "failed"
)

type TransactionWrapper struct {
	ledgerSequence uint32
	tx             ingest.LedgerTransaction
	ops            []transactionOperationWrapper
}

// type TransactionResultMeta struct {
// 	Result            TransactionResultPair
// 	FeeProcessing     LedgerEntryChanges
// 	TxApplyProcessing TransactionMeta
// }

func NewTransactionWrapper(tx ingest.LedgerTransaction, seq uint32) TransactionWrapper {
	var ops []transactionOperationWrapper
	for opi, op := range tx.Envelope.Operations() {
		operation := transactionOperationWrapper{
			index:          uint32(opi),
			transaction:    tx,
			operation:      op,
			ledgerSequence: seq,
		}

		ops = append(ops, operation)
	}

	return TransactionWrapper{
		ledgerSequence: seq,
		tx:             tx,
		ops:            ops,
	}
}

func (tw TransactionWrapper) GetTransactionHash() string {
	return tw.tx.Result.TransactionHash.HexString()
}

func (tw TransactionWrapper) GetStatus() string {
	if tw.tx.Result.Successful() {
		return SUCCESS
	}

	return FAILED
}

func (tw TransactionWrapper) GetLedgerSequence() uint32 {
	return tw.ledgerSequence
}

func (tw TransactionWrapper) GetApplicationOrder() uint32 {
	return tw.tx.Index
}

func (tw TransactionWrapper) GetEnvelopeXdr() string {
	bz, _ := tw.tx.Envelope.MarshalBinary()
	return string(bz)
}

func (tw TransactionWrapper) GetResultXdr() string {
	bz, _ := tw.tx.Result.MarshalBinary()
	return string(bz)
}

func (tw TransactionWrapper) GetResultMetaXdr() string {
	txResultMeta := xdr.TransactionResultMeta{
		Result:            tw.tx.Result,
		FeeProcessing:     tw.tx.FeeChanges,
		TxApplyProcessing: tw.tx.UnsafeMeta,
	}

	bz, _ := txResultMeta.MarshalBinary()

	return string(bz)
}

func (tw TransactionWrapper) GetModelsTransaction() *models.Transaction {
	return &models.Transaction{
		Hash:             tw.GetTransactionHash(),
		Status:           tw.GetStatus(),
		Ledger:           tw.GetLedgerSequence(),
		ApplicationOrder: tw.GetApplicationOrder(),
		EnvelopeXdr:      tw.GetEnvelopeXdr(),
		ResultXdr:        tw.GetResultXdr(),
		ResultMetaXdr:    tw.GetResultMetaXdr(),
	}
}
