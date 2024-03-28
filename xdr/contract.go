package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

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

	return result, errors.Errorf("Invalid contract id preimage type %v", f.Type)
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
		return result, nil
	}

	result.Address = address
	result.Salt = p.Salt.String()

	return result, nil
}
