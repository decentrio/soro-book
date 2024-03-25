package aggregation

import (
	"github.com/stellar/go/ingest"
)

type TransactionWrapper struct {
	tx  ingest.LedgerTransaction
	ops []transactionOperationWrapper
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
		tx:  tx,
		ops: ops,
	}
}
