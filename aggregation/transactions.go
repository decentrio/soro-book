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
	LedgerSequence uint32
	Tx             ingest.LedgerTransaction
	Ops            []transactionOperationWrapper
}

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
		EnvelopeXdr:      tw.GetEnvelopeXdr(),
		ResultXdr:        tw.GetResultXdr(),
		ResultMetaXdr:    tw.GetResultMetaXdr(),
		SourceAddress:    tw.Tx.Envelope.SourceAccount().ToAccountId().Address(),
	}
}
