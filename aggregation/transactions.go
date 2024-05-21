package aggregation

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/decentrio/converter/converter"
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
		if as.state == CONTRACT {
			if len(as.assetContractEventsQueue) == 0 && len(as.wasmContractEventsQueue) == 0 && len(as.contractDataEntrysQueue) == 0 {
				as.state = TX
			}
		}

		if as.state != TX {
			continue
		}

		if len(as.txQueue) == 0 {
			as.state = LEDGER
		}

		select {
		// Receive a new tx
		case tx := <-as.txQueue:
			as.handleReceiveNewTransaction(tx)

			as.state = CONTRACT
		// Terminate process
		case <-as.BaseService.Terminate():
			return
		default:
		}
		time.Sleep(time.Millisecond)
	}
}

func (as *Aggregation) handleReceiveNewTransaction(tw TransactionWrapper) {
	// tx := tw.GetModelsTransaction()
	// _, err := as.db.CreateTransaction(tx)
	// if err != nil {
	// 	as.Logger.Error(fmt.Sprintf("Error create ledger %d tx %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), err.Error()))
	// }

	evl, _ := converter.ConvertTransactionEnvelope(tw.Tx.Envelope)
	bz, _ := json.Marshal(evl)

	fmt.Printf("envelop: %s\n\n", bz)

	// Contract entry
	entries := tw.GetModelsContractDataEntry()
	for _, entry := range entries {
		// as.contractDataEntrysQueue <- entry
		bz, _ := json.Marshal(entry)
		fmt.Printf("entry: %s\n\n", bz)
	}

	wasmEvent, assetEvent, err := tw.GetContractEvents()
	if err != nil {
		return
	}
	// Soroban stellar asset events
	for _, e := range assetEvent {
		// as.assetContractEventsQueue <- e
		bz, _ := json.Marshal(e)
		fmt.Printf("assetEvent: %s\n\n", bz)
	}
	// Soroban wasm contract events
	for _, e := range wasmEvent {
		// as.wasmContractEventsQueue <- e
		bz, _ := json.Marshal(e)
		fmt.Printf("wasmEvent: %s\n", bz)
	}
}

type TransactionWrapper struct {
	LedgerSequence uint32
	Tx             ingest.LedgerTransaction
	Ops            []transactionOperationWrapper
	ProcessedAt    uint64
}

func NewTransactionWrapper(tx ingest.LedgerTransaction, ledgerSeq uint32, processedUnixTime uint64) TransactionWrapper {
	var ops []transactionOperationWrapper
	for opi, op := range tx.Envelope.Operations() {
		operation := transactionOperationWrapper{
			index:          uint32(opi),
			txIndex:        tx.Index,
			operation:      op,
			ledgerSequence: ledgerSeq,
		}

		ops = append(ops, operation)
	}

	return TransactionWrapper{
		LedgerSequence: ledgerSeq,
		Tx:             tx,
		Ops:            ops,
		ProcessedAt:    processedUnixTime,
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
		ProcessedAt:      tw.ProcessedAt,
	}
}
