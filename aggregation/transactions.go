package aggregation

import (
	"fmt"
	"time"

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
		time.Sleep(time.Millisecond)
	}
}

func (as *Aggregation) handleReceiveNewTransaction(tw TransactionWrapper) {
	tx := tw.GetModelsTransaction()
	_, err := as.db.CreateTransaction(tx)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error create ledger %d tx %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), err.Error()))
	}

	// Contract entry
	entries := tw.GetModelsContractDataEntry()
	for _, entry := range entries {
		go func(e models.Contract) {
			as.contractDataEntrysQueue <- e
		}(entry)
	}

	wasmEvent, assetEvent, err := tw.GetContractEvents()
	if err != nil {
		return
	}
	// Soroban stellar asset events
	for _, e := range assetEvent {
		go func(ae models.StellarAssetContractEvent) {
			as.assetContractEventsQueue <- ae
		}(e)
	}
	// Soroban wasm contract events
	for _, e := range wasmEvent {
		go func(we models.WasmContractEvent) {
			as.wasmContractEventsQueue <- we
		}(e)
	}
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
