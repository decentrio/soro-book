package converter

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

// func ConvertTransactionJSON(tx models.Transaction) models.TransactionJSON {
// 	return models.TransactionJSON{}
// }

func ConvertTransactionResultMeta(r xdr.TransactionResultMeta) (TransactionResultMeta, error) {
	var result TransactionResultMeta

	rs, err := ConvertTransactionResultPair(r.Result)
	if err != nil {
		return result, err
	}

	var fees LedgerEntryChanges
	for _, xdrFee := range r.FeeProcessing {
		fee, err := ConvertLedgerEntryChange(xdrFee)
		if err != nil {
			return result, err
		}
		fees = append(fees, fee)
	}

	txMeta, err := ConvertTransactionMeta(r.TxApplyProcessing)
	if err != nil {
		return result, err
	}

	result.Result = rs
	result.FeeProcessing = fees
	result.TxApplyProcessing = txMeta

	return result, nil
}

func ConvertTransactionMeta(m xdr.TransactionMeta) (TransactionMeta, error) {
	var result TransactionMeta

	switch m.V {
	case 0:
		var ops []OperationMeta
		for _, xdrOp := range *m.Operations {
			op, err := ConvertOperationMeta(xdrOp)
			if err != nil {
				return result, err
			}
			ops = append(ops, op)
		}

		result.Operations = &ops

		return result, nil
	case 1:
		v1, err := ConvertTransactionMetaV1(*m.V1)
		if err != nil {
			return result, err
		}
		result.V1 = &v1
		return result, nil
	case 2:
		v2, err := ConvertTransactionMetaV2(*m.V2)
		if err != nil {
			return result, err
		}
		result.V2 = &v2
		return result, nil
	case 3:
		v3, err := ConvertTransactionMetaV3(*m.V3)
		if err != nil {
			return result, err
		}
		result.V3 = &v3
		return result, nil
	}
	return result, errors.Errorf("error invalid TransactionMeta type %v", m.V)
}

func ConvertTransactionMetaV1(m xdr.TransactionMetaV1) (TransactionMetaV1, error) {
	var result TransactionMetaV1

	var txChanges LedgerEntryChanges
	for _, xdrTxChange := range m.TxChanges {
		txChange, err := ConvertLedgerEntryChange(xdrTxChange)
		if err != nil {
			return result, err
		}

		txChanges = append(txChanges, txChange)
	}

	var operations []OperationMeta
	for _, xdrOp := range m.Operations {
		op, err := ConvertOperationMeta(xdrOp)
		if err != nil {
			return result, err
		}

		operations = append(operations, op)
	}

	result.TxChanges = txChanges
	result.Operations = operations

	return result, nil
}

func ConvertTransactionMetaV2(m xdr.TransactionMetaV2) (TransactionMetaV2, error) {
	var result TransactionMetaV2

	var txChangesBefore LedgerEntryChanges
	for _, xdrTxChange := range m.TxChangesBefore {
		txChange, err := ConvertLedgerEntryChange(xdrTxChange)
		if err != nil {
			return result, err
		}

		txChangesBefore = append(txChangesBefore, txChange)
	}

	var operations []OperationMeta
	for _, xdrOp := range m.Operations {
		op, err := ConvertOperationMeta(xdrOp)
		if err != nil {
			return result, err
		}

		operations = append(operations, op)
	}

	var txChangesAfter LedgerEntryChanges
	for _, xdrTxChange := range m.TxChangesAfter {
		txChange, err := ConvertLedgerEntryChange(xdrTxChange)
		if err != nil {
			return result, err
		}

		txChangesAfter = append(txChangesAfter, txChange)
	}

	result.TxChangesBefore = txChangesBefore
	result.Operations = operations
	result.TxChangesAfter = txChangesAfter

	return result, nil
}

func ConvertTransactionMetaV3(m xdr.TransactionMetaV3) (TransactionMetaV3, error) {
	var result TransactionMetaV3

	ext := ConvertExtensionPoint(m.Ext)

	var txChangesBefore LedgerEntryChanges
	for _, xdrTxChange := range m.TxChangesBefore {
		txChange, err := ConvertLedgerEntryChange(xdrTxChange)
		if err != nil {
			return result, err
		}

		txChangesBefore = append(txChangesBefore, txChange)
	}

	var operations []OperationMeta
	for _, xdrOp := range m.Operations {
		op, err := ConvertOperationMeta(xdrOp)
		if err != nil {
			return result, err
		}

		operations = append(operations, op)
	}

	var txChangesAfter LedgerEntryChanges
	for _, xdrTxChange := range m.TxChangesAfter {
		txChange, err := ConvertLedgerEntryChange(xdrTxChange)
		if err != nil {
			return result, err
		}

		txChangesAfter = append(txChangesAfter, txChange)
	}

	var sorobanMeta SorobanTransactionMeta
	if m.SorobanMeta != nil {
		var err error
		sorobanMeta, err = ConvertSorobanTransactionMeta(*m.SorobanMeta)
		if err != nil {
			return result, err
		}
	}

	result.Ext = ext
	result.TxChangesBefore = txChangesBefore
	result.Operations = operations
	result.TxChangesAfter = txChangesAfter
	result.SorobanMeta = &sorobanMeta

	return result, nil
}

func ConvertSorobanTransactionMeta(m xdr.SorobanTransactionMeta) (SorobanTransactionMeta, error) {
	var result SorobanTransactionMeta
	ext := ConvertExtensionPoint(m.Ext)

	var events []ContractEvent
	for _, xdrEvent := range m.Events {
		event, err := ConvertContractEvent(xdrEvent)
		if err != nil {
			return result, err
		}
		events = append(events, event)
	}

	returnValue, err := ConvertScVal(m.ReturnValue)
	if err != nil {
		return result, err
	}

	var diagnosticEvents []DiagnosticEvent
	for _, xdrEvent := range m.DiagnosticEvents {
		event, err := ConvertDiagnosticEvent(xdrEvent)
		if err != nil {
			return result, err
		}
		diagnosticEvents = append(diagnosticEvents, event)
	}

	result.Ext = ext
	result.Events = events
	result.ReturnValue = returnValue
	result.DiagnosticEvents = diagnosticEvents

	return result, nil
}

func ConvertTransactionResultPair(r xdr.TransactionResultPair) (TransactionResultPair, error) {
	var result TransactionResultPair
	result.TransactionHash = r.TransactionHash.HexString()

	rs, err := ConvertTransactionResult(r.Result)
	if err != nil {
		return result, err
	}
	result.Result = rs

	return result, nil
}

func ConvertTransactionResult(r xdr.TransactionResult) (TransactionResult, error) {
	var result TransactionResult
	result.FeeCharged = int64(r.FeeCharged)

	rs, err := ConvertTransactionResultResult(r.Result)
	if err != nil {
		return result, err
	}
	result.Result = rs

	ext := ConvertTransactionResultExt(r.Ext)
	result.Ext = ext

	return result, nil
}

func ConvertTransactionResultResult(r xdr.TransactionResultResult) (TransactionResultResult, error) {
	var result TransactionResultResult
	result.Code = int32(r.Code)

	if r.Code == xdr.TransactionResultCodeTxFeeBumpInnerSuccess || r.Code == xdr.TransactionResultCodeTxFeeBumpInnerFailed {
		innerResult, err := ConvertInnerTransactionResultPair(*r.InnerResultPair)
		if err != nil {
			return result, err
		}
		result.InnerResultPair = &innerResult
	} else if r.Code == xdr.TransactionResultCodeTxSuccess || r.Code == xdr.TransactionResultCodeTxFailed {
		var opResult []OperationResult
		if r.Results != nil {
			for _, xdrResult := range *r.Results {
				op, err := ConvertOperationResult(xdrResult)
				if err != nil {
					return result, err
				}

				opResult = append(opResult, op)
			}
			result.Results = &opResult
		}
	}
	return result, nil
}

func ConvertInnerTransactionResultPair(r xdr.InnerTransactionResultPair) (InnerTransactionResultPair, error) {
	var result InnerTransactionResultPair
	result.TransactionHash = r.TransactionHash.HexString()

	rs, err := ConvertInnerTransactionResult(r.Result)
	if err != nil {
		return result, err
	}
	result.Result = rs

	return result, nil
}

func ConvertInnerTransactionResult(r xdr.InnerTransactionResult) (InnerTransactionResult, error) {
	var result InnerTransactionResult
	result.FeeCharged = int64(r.FeeCharged)

	rs, err := ConvertInnerTransactionResultResult(r.Result)
	if err != nil {
		return result, err
	}
	result.Result = rs

	ext := ConvertInnerTransactionResultExt(r.Ext)
	result.Ext = ext

	return result, nil
}

func ConvertInnerTransactionResultResult(r xdr.InnerTransactionResultResult) (InnerTransactionResultResult, error) {
	var result InnerTransactionResultResult
	result.Code = int32(r.Code)

	if r.Code == xdr.TransactionResultCodeTxSuccess || r.Code == xdr.TransactionResultCodeTxFailed {
		var opResults []OperationResult
		for _, xdrResult := range *r.Results {
			r, err := ConvertOperationResult(xdrResult)
			if err != nil {
				return result, err
			}
			opResults = append(opResults, r)
		}

		result.Results = &opResults
	}

	return result, nil
}

func ConvertInnerTransactionResultExt(e xdr.InnerTransactionResultExt) InnerTransactionResultExt {
	return InnerTransactionResultExt{V: e.V}
}

func ConvertTransactionResultExt(e xdr.TransactionResultExt) TransactionResultExt {
	return TransactionResultExt{V: e.V}
}

func ConvertFeeBumpTransaction(tx xdr.FeeBumpTransaction) (FeeBumpTransaction, error) {
	var result FeeBumpTransaction

	feeSource, err := ConvertMuxedAccount(tx.FeeSource)
	if err != nil {
		return result, err
	}

	innerTx, err := ConvertFeeBumpTransactionInnerTx(tx.InnerTx)
	if err != nil {
		return result, err
	}

	ext := ConvertFeeBumpTransactionExt(tx.Ext)

	result.FeeSource = feeSource
	result.Fee = int64(tx.Fee)
	result.InnerTx = innerTx
	result.Ext = ext

	return result, nil
}

func ConvertTransaction(tx xdr.Transaction) (Transaction, error) {
	var result Transaction

	sourceAccount, err := ConvertMuxedAccount(tx.SourceAccount)
	if err != nil {
		return result, err
	}

	cond, err := ConvertPreconditions(tx.Cond)
	if err != nil {
		return result, err
	}

	memo, err := ConvertMemo(tx.Memo)
	if err != nil {
		return result, err
	}

	var ops []Operation
	for _, xdrOp := range tx.Operations {
		op, err := ConvertOperation(xdrOp)
		if err != nil {
			return result, err
		}
		ops = append(ops, op)
	}

	ext, err := ConvertTxExt(tx.Ext)
	if err != nil {
		return result, err
	}

	result.SourceAccount = sourceAccount
	result.Fee = uint32(tx.Fee)
	result.SeqNum = int64(tx.SeqNum)
	result.Cond = cond
	result.Memo = memo
	result.Operations = ops
	result.Ext = ext

	return result, nil

}

// TODO: testing
func ConvertTransactionV0(tx xdr.TransactionV0) (TransactionV0, error) {
	var txV0 TransactionV0

	txV0.SourceAccountEd25519 = tx.SourceAccountEd25519.String()
	txV0.Fee = uint32(tx.Fee)
	txV0.SeqNum = int64(tx.SeqNum)

	tb, err := ConvertTimeBounds(tx.TimeBounds)
	if err != nil {
		return txV0, err
	}
	txV0.TimeBounds = tb

	memo, err := ConvertMemo(tx.Memo)
	if err != nil {
		return txV0, err
	}
	txV0.Memo = memo

	var ops []Operation
	for _, opXdr := range tx.Operations {
		op, err := ConvertOperation(opXdr)
		if err != nil {
			return txV0, err
		}

		ops = append(ops, op)
	}
	txV0.Operations = ops

	ext, err := ConvertTxV0Ext(tx.Ext)
	if err != nil {
		return txV0, err
	}
	txV0.Ext = ext

	return txV0, nil
}

// TODO: testing
func ConvertTimeBounds(tb *xdr.TimeBounds) (*TimeBounds, error) {
	if tb == nil {
		return nil, nil
	}

	return &TimeBounds{
		MinTime: uint64(tb.MinTime),
		MaxTime: uint64(tb.MaxTime),
	}, nil
}

// TODO: testing
func ConvertMemo(memo xdr.Memo) (Memo, error) {
	var result Memo

	switch memo.Type {
	case xdr.MemoTypeMemoNone:
		return result, nil
	case xdr.MemoTypeMemoText:
		text, ok := memo.GetText()
		if !ok {
			return result, errors.Errorf("error invalid memo type text %v", memo)
		}
		result.Text = &text

		return result, nil
	case xdr.MemoTypeMemoId:
		xdrId, ok := memo.GetId()
		if !ok {
			return result, errors.Errorf("error invalid memo type id %v", memo)
		}
		id := uint64(xdrId)
		result.Id = &id

		return result, nil
	case xdr.MemoTypeMemoHash:
		xdrHash, ok := memo.GetHash()
		if !ok {
			return result, errors.Errorf("error invalid memo type hash %v", memo)
		}
		hash := xdrHash.HexString()
		result.Hash = &hash

		return result, nil
	case xdr.MemoTypeMemoReturn:
		xdrRetHash, ok := memo.GetRetHash()
		if !ok {
			return result, errors.Errorf("error invalid memo type rethash%v", memo)
		}
		retHash := xdrRetHash.HexString()
		result.RetHash = &retHash

		return result, nil
	default:
		return result, errors.Errorf("error invalid memo %v", memo)
	}
}

// TODO: testing
func ConvertTxV0Ext(e xdr.TransactionV0Ext) (TransactionV0Ext, error) {
	return TransactionV0Ext{
		V: e.V,
	}, nil
}

func ConvertTxExt(e xdr.TransactionExt) (TransactionExt, error) {
	var result TransactionExt

	var data SorobanTransactionData
	var err error
	if e.SorobanData != nil {
		data, err = ConvertSorobanTransactionData(*e.SorobanData)
		if err != nil {
			return result, err
		}
	}

	result.V = e.V
	result.SorobanData = &data

	return result, nil
}

func ConvertFeeBumpTransactionInnerTx(f xdr.FeeBumpTransactionInnerTx) (FeeBumpTransactionInnerTx, error) {
	var result FeeBumpTransactionInnerTx
	switch f.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTx:
		v1, err := ConvertTransactionV1Envelope(f.V1)
		if err != nil {
			return result, err
		}
		result.V1 = &v1

		return result, nil
	}
	return result, errors.Errorf("error invalid FeeBumpTransactionInnerTx %v", f.Type)
}

func ConvertFeeBumpTransactionExt(f xdr.FeeBumpTransactionExt) FeeBumpTransactionExt {
	return FeeBumpTransactionExt{V: f.V}
}

func ConvertDataEntry(e xdr.DataEntry) DataEntry {
	accountId := PublicKey{
		Ed25519: e.AccountId.Ed25519.String(),
	}

	return DataEntry{
		AccountId: accountId,
		DataName:  string(e.DataName),
		DataValue: e.DataValue,
		Ext:       ConvertDataEntryExt(e.Ext),
	}
}

func ConvertDataEntryExt(e xdr.DataEntryExt) DataEntryExt {
	return DataEntryExt{V: e.V}
}

func ConvertTtlEntry(e xdr.TtlEntry) TtlEntry {
	return TtlEntry{
		KeyHash:            e.KeyHash.HexString(),
		LiveUntilLedgerSeq: uint32(e.LiveUntilLedgerSeq),
	}
}

func ConvertConfigSettingEntry(e xdr.ConfigSettingEntry) (ConfigSettingEntry, error) {
	var result ConfigSettingEntry

	result.ConfigSettingId = int32(e.ConfigSettingId)

	switch e.ConfigSettingId {
	case xdr.ConfigSettingIdConfigSettingContractMaxSizeBytes:
		value := uint32(*e.ContractMaxSizeBytes)
		result.ContractMaxSizeBytes = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractComputeV0:
		value := ConvertConfigSettingContractComputeV0(*e.ContractCompute)
		result.ContractCompute = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractLedgerCostV0:
		value := ConvertConfigSettingContractLedgerCostV0(*e.ContractLedgerCost)
		result.ContractLedgerCost = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractHistoricalDataV0:
		value := ConvertConfigSettingContractHistoricalDataV0(*e.ContractHistoricalData)
		result.ContractHistoricalData = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractEventsV0:
		value := ConvertConfigSettingContractEventsV0(*e.ContractEvents)
		result.ContractEvents = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractBandwidthV0:
		value := ConvertConfigSettingContractBandwidthV0(*e.ContractBandwidth)
		result.ContractBandwidth = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractCostParamsCpuInstructions:
		value := ConvertContractCostParams(*e.ContractCostParamsCpuInsns)
		result.ContractCostParamsCpuInsns = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractCostParamsMemoryBytes:
		value := ConvertContractCostParams(*e.ContractCostParamsMemBytes)
		result.ContractCostParamsMemBytes = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractDataKeySizeBytes:
		value := uint32(*e.ContractDataKeySizeBytes)
		result.ContractDataKeySizeBytes = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractDataEntrySizeBytes:
		value := uint32(*e.ContractDataEntrySizeBytes)
		result.ContractDataEntrySizeBytes = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingStateArchival:
		value := ConvertStateArchivalSettings(*e.StateArchivalSettings)
		result.StateArchivalSettings = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingContractExecutionLanes:
		value := ConvertConfigSettingContractExecutionLanesV0(*e.ContractExecutionLanes)
		result.ContractExecutionLanes = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingBucketlistSizeWindow:
		var value []uint64
		for _, s := range *e.BucketListSizeWindow {
			value = append(value, uint64(s))
		}

		result.BucketListSizeWindow = &value
		return result, nil
	case xdr.ConfigSettingIdConfigSettingEvictionIterator:
		value := ConvertEvictionIterator(*e.EvictionIterator)
		result.EvictionIterator = &value
		return result, nil
	}
	return result, errors.Errorf("invalid ConfigSettingEntry code id %v", e.ConfigSettingId)
}
