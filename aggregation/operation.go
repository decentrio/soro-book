package aggregation

import (
	"fmt"

	"github.com/decentrio/soro-book/database/models"
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

func (operation *transactionOperationWrapper) GetContractEvents() map[models.ContractEvent][]string {
	var eventsMap = make(map[models.ContractEvent][]string)
	var order = uint32(1)

	for _, event := range operation.transaction.UnsafeMeta.V3.SorobanMeta.Events {
		var topics []string
		eventType, found := getEventType(event.Body)
		if !found {
			continue
		}

		for _, topic := range event.Body.V0.Topics {
			bz, err := topic.MarshalBinary()
			if err != nil {
				break
			}

			topics = append(topics, string(bz))
		}

		valueBz, err := event.Body.V0.Data.MarshalBinary()
		if err != nil {
			continue
		}

		event := models.ContractEvent{
			Id:         fmt.Sprintf("%019d-%010d", operation.ID(), order), // ID should be combine from operation ID and event index
			ContractId: event.ContractId.HexString(),
			LedgerSeq:  operation.ledgerSequence,
			TxHash:     operation.transaction.Result.TransactionHash.HexString(),
			EventType:  eventType,
			Value:      string(valueBz),
		}

		eventsMap[event] = topics

		order++
	}

	return eventsMap
}
