package aggregation

import (
	"github.com/stellar/go/ingest"
	"github.com/stellar/go/toid"
	"github.com/stellar/go/xdr"
)

type transactionOperationWrapper struct {
	index          uint32
	transaction    ingest.LedgerTransaction
	operation      xdr.Operation
	ledgerSequence uint32
}

// ID returns the ID for the operation.
func (operation *transactionOperationWrapper) ID() int64 {
	return toid.New(
		int32(operation.ledgerSequence),
		int32(operation.transaction.Index),
		int32(operation.index+1),
	).ToInt64()
}

// TransactionID returns the id for the transaction related with this operation.
func (operation *transactionOperationWrapper) TransactionID() int64 {
	return toid.New(int32(operation.ledgerSequence), int32(operation.transaction.Index), 0).ToInt64()
}

// SourceAccount returns the operation's source account.
func (operation *transactionOperationWrapper) SourceAccount() *xdr.MuxedAccount {
	sourceAccount := operation.operation.SourceAccount
	if sourceAccount != nil {
		return sourceAccount
	} else {
		ret := operation.transaction.Envelope.SourceAccount()
		return &ret
	}
}

// OperationType returns the operation type.
func (operation *transactionOperationWrapper) OperationType() xdr.OperationType {
	return operation.operation.Body.Type
}
