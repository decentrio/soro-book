package aggregation

import (
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
	tx := tw.GetModelsTransaction()
	_, err := as.db.CreateTransaction(tx)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("error create ledger %d tx %s: %s", tw.GetLedgerSequence(), tw.GetTransactionHash(), err.Error()))
	}

	// if this is invokeHostFuncTx, we should store the detail
	invokeHostFuncTx, createContractTx, err := isInvokeHostFunctionTx(tw.Tx, tw.LedgerSequence)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("error invoke host function %s", err.Error()))
	}

	for _, ivhft := range invokeHostFuncTx {
		_, err := as.db.CreateContractInvokedTransaction(&ivhft)
		if err != nil {
			as.Logger.Error(fmt.Sprintf("error create invoke host function %s", err.Error()))
		}
	}

	for _, cct := range createContractTx {
		_, err := as.db.CreateContractCreatedTransaction(&cct)
		if err != nil {
			as.Logger.Error(fmt.Sprintf("error create contract created function %s", err.Error()))
		}
	}

	// Contract entry
	entries := tw.GetModelsContractDataEntry()
	for _, entry := range entries {
		as.contractDataEntrysQueue <- entry
	}

	wasmEvent, assetEvent, err := tw.GetContractEvents()
	if err != nil {
		return
	}
	// Soroban stellar asset events
	for _, e := range assetEvent {
		as.assetContractEventsQueue <- e
	}
	// Soroban wasm contract events
	for _, e := range wasmEvent {
		as.wasmContractEventsQueue <- e
	}
}

func isInvokeHostFunctionTx(tx ingest.LedgerTransaction, ledgerSeq uint32) ([]models.InvokeTransaction, []models.ContractsCode, error) {
	var invokeFuncTxs []models.InvokeTransaction
	var createdContracts []models.ContractsCode

	ops := tx.Envelope.Operations()
	for _, op := range ops {
		if op.Body.Type == xdr.OperationTypeInvokeHostFunction {
			ihfOp := op.Body.MustInvokeHostFunctionOp()
			switch ihfOp.HostFunction.Type {
			case xdr.HostFunctionTypeHostFunctionTypeInvokeContract:

				ic := ihfOp.HostFunction.MustInvokeContract()
				ca, err := converter.ConvertScAddress(ic.ContractAddress)
				if err != nil {
					continue
				}

				fn := string(ic.FunctionName)

				args, err := ic.MarshalBinary()
				if err != nil {
					continue
				}

				var invokeFuncTx models.InvokeTransaction
				invokeFuncTx.Hash = tx.Result.TransactionHash.HexString()
				invokeFuncTx.ContractId = *ca.ContractId
				invokeFuncTx.FunctionType = "invoke_host_function"
				invokeFuncTx.FunctionName = fn
				invokeFuncTx.Args = args

				invokeFuncTxs = append(invokeFuncTxs, invokeFuncTx)

			case xdr.HostFunctionTypeHostFunctionTypeCreateContract:
				ccop := ihfOp.HostFunction.MustCreateContract()

				var createContractTx models.ContractsCode
				creator := tx.Envelope.SourceAccount().ToAccountId().Address()

				contractId, found := getCreatedContractId(tx.Envelope)
				if !found {
					continue
				}

				var contractCode string
				if ccop.Executable.WasmHash != nil {
					contractCode = (*ccop.Executable.WasmHash).HexString()
				}

				createContractTx.CreatorAddress = creator
				createContractTx.ContractId = contractId
				createContractTx.ContractCode = contractCode
				createContractTx.CreatedLedger = ledgerSeq

				createdContracts = append(createdContracts, createContractTx)

			case xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm:
				// we do not care about this type
				continue
			}

		}
	}

	return invokeFuncTxs, createdContracts, nil
}

func getCreatedContractId(op xdr.TransactionEnvelope) (string, bool) {

	switch op.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
		return "", false
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		v1 := op.MustV1()
		ext := v1.Tx.Ext
		sorobanData := ext.MustSorobanData()

		footprints := sorobanData.Resources.Footprint.ReadWrite
		for _, fp := range footprints {
			if fp.Type == xdr.LedgerEntryTypeContractData {
				contractData := fp.MustContractData()
				contractId, _ := converter.ConvertScAddress(contractData.Contract)
				return *contractId.ContractId, true
			}
		}

		return "", false
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		return "", false
	default:
		return "", false
	}
}

type TransactionWrapper struct {
	LedgerSequence uint32
	Tx             ingest.LedgerTransaction
	Ops            []transactionOperationWrapper
	Time           uint64
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
		Time:           processedUnixTime,
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
		TransactionTime:  tw.Time,
	}
}
