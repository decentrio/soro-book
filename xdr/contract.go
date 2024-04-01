package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertSorobanAuthorizationEntry(e xdr.SorobanAuthorizationEntry) (SorobanAuthorizationEntry, error) {
	var result SorobanAuthorizationEntry

	credentials, err := ConvertSorobanCredentials(e.Credentials)
	if err != nil {
		return result, err
	}

	rootInvocation, err := ConvertSorobanAuthorizedInvocation(e.RootInvocation)
	if err != nil {
		return result, err
	}

	result.Credentials = credentials
	result.RootInvocation = rootInvocation

	return result, nil
}

func ConvertSorobanCredentials(c xdr.SorobanCredentials) (SorobanCredentials, error) {
	var result SorobanCredentials
	switch c.Type {
	case xdr.SorobanCredentialsTypeSorobanCredentialsSourceAccount:
		// void
		return result, nil
	case xdr.SorobanCredentialsTypeSorobanCredentialsAddress:
		address, err := ConvertSorobanAddressCredentials(*c.Address)
		if err != nil {
			return result, err
		}

		result.Address = &address
		return result, nil
	}

	return result, errors.Errorf("Invalid ConvertSorobanCredentials type %v\n", c.Type)
}

func ConvertSorobanAddressCredentials(c xdr.SorobanAddressCredentials) (SorobanAddressCredentials, error) {
	var result SorobanAddressCredentials

	address, err := ConvertScAddress(c.Address)
	if err != nil {
		return result, err
	}

	signature, err := ConvertScVal(c.Signature)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Nonce = int64(c.Nonce)
	result.SignatureExpirationLedger = uint32(c.SignatureExpirationLedger)
	result.Signature = signature

	return result, nil
}

func ConvertSorobanAuthorizedInvocation(i xdr.SorobanAuthorizedInvocation) (SorobanAuthorizedInvocation, error) {
	var result SorobanAuthorizedInvocation
	function, err := ConvertSorobanAuthorizedFunction(i.Function)
	if err != nil {
		return result, err
	}
	result.Function = function

	var subs []SorobanAuthorizedInvocation
	for _, xdrSub := range i.SubInvocations {
		sub, err := ConvertSorobanAuthorizedInvocation(xdrSub)
		if err != nil {
			return result, err
		}

		subs = append(subs, sub)
	}
	result.SubInvocations = subs

	return result, nil
}

func ConvertSorobanAuthorizedFunction(f xdr.SorobanAuthorizedFunction) (SorobanAuthorizedFunction, error) {
	var result SorobanAuthorizedFunction
	switch f.Type {
	case xdr.SorobanAuthorizedFunctionTypeSorobanAuthorizedFunctionTypeContractFn:
		contractFn, err := ConvertInvokeContractArgs(*f.ContractFn)
		if err != nil {
			return result, err
		}
		result.ContractFn = &contractFn

		return result, nil
	case xdr.SorobanAuthorizedFunctionTypeSorobanAuthorizedFunctionTypeCreateContractHostFn:
		createContract, err := ConvertCreateContractArgs(*f.CreateContractHostFn)
		if err != nil {
			return result, err
		}
		result.CreateContractHostFn = &createContract

		return result, nil
	}

	return result, errors.Errorf("Invalid SorobanAuthorizedFunction type %v", f.Type)
}

func ConvertSorobanTransactionData(d xdr.SorobanTransactionData) (SorobanTransactionData, error) {
	var result SorobanTransactionData

	resources, err := ConvertSorobanResources(d.Resources)
	if err != nil {
		return result, err
	}

	result.Ext = ConvertExtensionPoint(d.Ext)
	result.Resources = resources
	result.ResourceFee = int64(d.ResourceFee)

	return result, nil
}

func ConvertSorobanResources(r xdr.SorobanResources) (SorobanResources, error) {
	var result SorobanResources

	footPrint, err := ConvertLedgerFootprint(r.Footprint)
	if err != nil {
		return result, err
	}

	result.Footprint = footPrint
	result.Instructions = uint32(r.Instructions)
	result.ReadBytes = result.ReadBytes
	result.WriteBytes = result.WriteBytes

	return result, nil
}

func ConvertHostFunction(f xdr.HostFunction) (HostFunction, error) {
	var result HostFunction
	switch f.Type {
	case xdr.HostFunctionTypeHostFunctionTypeInvokeContract:
		invokeContract, err := ConvertInvokeContractArgs(*f.InvokeContract)
		if err != nil {
			return result, err
		}
		result.InvokeContract = &invokeContract

		return result, nil
	case xdr.HostFunctionTypeHostFunctionTypeCreateContract:
		createContract, err := ConvertCreateContractArgs(*f.CreateContract)
		if err != nil {
			return result, err
		}
		result.CreateContract = &createContract

		return result, nil
	case xdr.HostFunctionTypeHostFunctionTypeUploadContractWasm:
		wasm := *f.Wasm
		result.Wasm = &wasm

		return result, nil
	}

	return result, errors.Errorf("Invalid host function type %v", f.Type)
}

func ConvertInvokeContractArgs(a xdr.InvokeContractArgs) (InvokeContractArgs, error) {
	var result InvokeContractArgs

	contractAddress, err := ConvertScAddress(a.ContractAddress)
	if err != nil {
		return result, err
	}

	funcName := ScSymbol(a.FunctionName)

	var args []ScVal
	for _, xdrArg := range a.Args {
		arg, err := ConvertScVal(xdrArg)
		if err != nil {
			return result, err
		}

		args = append(args, arg)
	}

	result.ContractAddress = contractAddress
	result.FunctionName = funcName
	result.Args = args

	return result, nil
}

func ConvertCreateContractArgs(a xdr.CreateContractArgs) (CreateContractArgs, error) {
	var result CreateContractArgs

	contractIdPreimage, err := ConvertContractIdPreimage(a.ContractIdPreimage)
	if err != nil {
		return result, err
	}

	executable, err := ConvertContractExecutable(a.Executable)
	if err != nil {
		return result, err
	}

	result.ContractIdPreimage = contractIdPreimage
	result.Executable = executable

	return result, nil
}

func ConvertContractExecutable(e xdr.ContractExecutable) (ContractExecutable, error) {
	var result ContractExecutable
	switch e.Type {
	case xdr.ContractExecutableTypeContractExecutableWasm:
		wasmHash := (*e.WasmHash).HexString()
		result.WasmHash = &wasmHash
		return result, nil
	case xdr.ContractExecutableTypeContractExecutableStellarAsset:
		return result, nil
	}

	return result, errors.Errorf("Invalid contract executable type %v", e.Type)
}

func ConvertContractIdPreimage(p xdr.ContractIdPreimage) (ContractIdPreimage, error) {
	var result ContractIdPreimage

	switch p.Type {
	case xdr.ContractIdPreimageTypeContractIdPreimageFromAddress:
		fromAddress, err := ConvertContractIdPreimageFromAddress(*p.FromAddress)
		if err != nil {
			return result, err
		}
		result.FromAddress = &fromAddress

		return result, nil
	case xdr.ContractIdPreimageTypeContractIdPreimageFromAsset:
		fromAsset, err := ConvertAsset(*p.FromAsset)
		if err != nil {
			return result, err
		}
		result.FromAsset = &fromAsset

		return result, nil
	}
	return result, errors.Errorf("Invalid contract id preimage type %v", p.Type)
}

func ConvertContractIdPreimageFromAddress(p xdr.ContractIdPreimageFromAddress) (ContractIdPreimageFromAddress, error) {
	var result ContractIdPreimageFromAddress
	address, err := ConvertScAddress(p.Address)
	if err != nil {
		return result, err
	}

	result.Address = address
	result.Salt = p.Salt.String()

	return result, nil
}

func ConvertContractCodeEntry(e xdr.ContractCodeEntry) ContractCodeEntry {
	return ContractCodeEntry{
		Ext:  ConvertExtensionPoint(e.Ext),
		Hash: e.Hash.HexString(),
		Code: e.Code,
	}
}

func ConvertContractDataEntry(e xdr.ContractDataEntry) (ContractDataEntry, error) {
	var result ContractDataEntry

	ext := ConvertExtensionPoint(e.Ext)

	contract, err := ConvertScAddress(e.Contract)
	if err != nil {
		return result, err
	}

	key, err := ConvertScVal(e.Key)
	if err != nil {
		return result, err
	}

	val, err := ConvertScVal(e.Val)
	if err != nil {
		return result, err
	}

	result.Ext = ext
	result.Contract = contract
	result.Key = key
	result.Durability = int32(e.Durability)
	result.Val = val

	return result, nil
}

func ConvertConfigSettingContractComputeV0(c xdr.ConfigSettingContractComputeV0) ConfigSettingContractComputeV0 {
	return ConfigSettingContractComputeV0{
		LedgerMaxInstructions:           int64(c.LedgerMaxInstructions),
		TxMaxInstructions:               int64(c.TxMaxInstructions),
		FeeRatePerInstructionsIncrement: int64(c.FeeRatePerInstructionsIncrement),
		TxMemoryLimit:                   uint32(c.TxMemoryLimit),
	}
}

func ConvertConfigSettingContractLedgerCostV0(c xdr.ConfigSettingContractLedgerCostV0) ConfigSettingContractLedgerCostV0 {
	return ConfigSettingContractLedgerCostV0{
		LedgerMaxReadLedgerEntries:     uint32(c.LedgerMaxReadLedgerEntries),
		LedgerMaxReadBytes:             uint32(c.LedgerMaxReadBytes),
		LedgerMaxWriteLedgerEntries:    uint32(c.LedgerMaxWriteLedgerEntries),
		LedgerMaxWriteBytes:            uint32(c.LedgerMaxWriteBytes),
		TxMaxReadLedgerEntries:         uint32(c.TxMaxReadLedgerEntries),
		TxMaxReadBytes:                 uint32(c.TxMaxReadBytes),
		TxMaxWriteLedgerEntries:        uint32(c.TxMaxWriteLedgerEntries),
		TxMaxWriteBytes:                uint32(c.TxMaxWriteBytes),
		FeeReadLedgerEntry:             int64(c.FeeReadLedgerEntry),
		FeeWriteLedgerEntry:            int64(c.FeeWriteLedgerEntry),
		FeeRead1Kb:                     int64(c.FeeRead1Kb),
		BucketListTargetSizeBytes:      int64(c.BucketListTargetSizeBytes),
		WriteFee1KbBucketListLow:       int64(c.WriteFee1KbBucketListLow),
		WriteFee1KbBucketListHigh:      int64(c.WriteFee1KbBucketListHigh),
		BucketListWriteFeeGrowthFactor: uint32(c.BucketListWriteFeeGrowthFactor),
	}
}

func ConvertConfigSettingContractHistoricalDataV0(c xdr.ConfigSettingContractHistoricalDataV0) ConfigSettingContractHistoricalDataV0 {
	return ConfigSettingContractHistoricalDataV0{FeeHistorical1Kb: int64(c.FeeHistorical1Kb)}
}

func ConvertConfigSettingContractEventsV0(c xdr.ConfigSettingContractEventsV0) ConfigSettingContractEventsV0 {
	return ConfigSettingContractEventsV0{
		TxMaxContractEventsSizeBytes: uint32(c.TxMaxContractEventsSizeBytes),
		FeeContractEvents1Kb:         int64(c.FeeContractEvents1Kb),
	}
}

func ConvertConfigSettingContractBandwidthV0(c xdr.ConfigSettingContractBandwidthV0) ConfigSettingContractBandwidthV0 {
	return ConfigSettingContractBandwidthV0{
		LedgerMaxTxsSizeBytes: uint32(c.LedgerMaxTxsSizeBytes),
		TxMaxSizeBytes:        uint32(c.TxMaxSizeBytes),
		FeeTxSize1Kb:          int64(c.FeeTxSize1Kb),
	}
}

func ConvertContractCostParams(c xdr.ContractCostParams) ContractCostParams {
	var result ContractCostParams
	for _, xdrEntry := range c {
		entry := ConvertContractCostParamEntry(xdrEntry)
		result = append(result, entry)
	}
	return result
}

func ConvertContractCostParamEntry(c xdr.ContractCostParamEntry) ContractCostParamEntry {
	return ContractCostParamEntry{
		Ext:        ConvertExtensionPoint(c.Ext),
		ConstTerm:  int64(c.ConstTerm),
		LinearTerm: int64(c.LinearTerm),
	}
}

func ConvertConfigSettingContractExecutionLanesV0(c xdr.ConfigSettingContractExecutionLanesV0) ConfigSettingContractExecutionLanesV0 {
	return ConfigSettingContractExecutionLanesV0{LedgerMaxTxCount: uint32(c.LedgerMaxTxCount)}
}

func ConvertEvictionIterator(i xdr.EvictionIterator) EvictionIterator {
	return EvictionIterator{
		BucketListLevel:  uint32(i.BucketListLevel),
		IsCurrBucket:     i.IsCurrBucket,
		BucketFileOffset: uint64(i.BucketFileOffset),
	}
}
