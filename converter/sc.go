package converter

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertScAddress(a xdr.ScAddress) (ScAddress, error) {
	var result ScAddress
	switch a.Type {
	case xdr.ScAddressTypeScAddressTypeAccount:
		account := PublicKey{
			Ed25519: ConvertEd25519(a.AccountId.Ed25519),
		}
		result.AccountId = &account

		return result, nil
	case xdr.ScAddressTypeScAddressTypeContract:
		contract := a.ContractId.HexString()
		result.ContractId = &contract

		return result, nil
	}
	return result, errors.Errorf("error invalid ScAddress type %v", a.Type)
}

func ConvertScError(e xdr.ScError) (ScError, error) {
	var result ScError
	switch e.Type {
	case xdr.ScErrorTypeSceContract:
		contractCode := uint32(*e.ContractCode)
		result.ContractCode = &contractCode
		return result, nil
	case xdr.ScErrorTypeSceWasmVm,
		xdr.ScErrorTypeSceContext,
		xdr.ScErrorTypeSceStorage,
		xdr.ScErrorTypeSceObject,
		xdr.ScErrorTypeSceCrypto,
		xdr.ScErrorTypeSceEvents,
		xdr.ScErrorTypeSceBudget,
		xdr.ScErrorTypeSceValue,
		xdr.ScErrorTypeSceAuth:
		code := int32(*e.Code)
		result.Code = &code
		return result, nil
	}
	return result, errors.Errorf("error invalid ScError type %v", e.Type)
}

func ConvertScVal(v xdr.ScVal) (ScVal, error) {
	var result ScVal
	switch v.Type {
	case xdr.ScValTypeScvBool:
		b := *v.B
		result.B = &b
		return result, nil
	case xdr.ScValTypeScvVoid:
		// void
		return result, nil
	case xdr.ScValTypeScvError:
		scErr, err := ConvertScError(*v.Error)
		if err != nil {
			return result, err
		}
		result.Error = &scErr

		return result, nil
	case xdr.ScValTypeScvU32:
		u32 := uint32(*v.U32)
		result.U32 = &u32
		return result, nil
	case xdr.ScValTypeScvI32:
		i32 := int32(*v.I32)
		result.I32 = &i32
		return result, nil
	case xdr.ScValTypeScvU64:
		u64 := uint64(*v.U64)
		result.U64 = &u64
		return result, nil
	case xdr.ScValTypeScvI64:
		i64 := int64(*v.I64)
		result.I64 = &i64
		return result, nil
	case xdr.ScValTypeScvTimepoint:
		tp := uint64(*v.Timepoint)
		result.Timepoint = &tp
		return result, nil
	case xdr.ScValTypeScvDuration:
		duration := uint64(*v.Duration)
		result.Duration = &duration
		return result, nil
	case xdr.ScValTypeScvU128:
		xdrU128 := *v.U128
		u128 := UInt128Parts{
			Hi: uint64(xdrU128.Hi),
			Lo: uint64(xdrU128.Lo),
		}
		result.U128 = &u128
		return result, nil
	case xdr.ScValTypeScvI128:
		xdrI128 := *v.I128
		i128 := Int128Parts{
			Hi: int64(xdrI128.Hi),
			Lo: uint64(xdrI128.Lo),
		}
		result.I128 = &i128
		return result, nil
	case xdr.ScValTypeScvU256:
		xdrU256 := *v.U256
		u256 := UInt256Parts{
			HiHi: uint64(xdrU256.HiHi),
			HiLo: uint64(xdrU256.HiLo),
			LoHi: uint64(xdrU256.LoHi),
			LoLo: uint64(xdrU256.LoLo),
		}
		result.U256 = &u256
		return result, nil
	case xdr.ScValTypeScvI256:
		xdrI256 := *v.I256
		i256 := Int256Parts{
			HiHi: int64(xdrI256.HiHi),
			HiLo: uint64(xdrI256.HiLo),
			LoHi: uint64(xdrI256.LoHi),
			LoLo: uint64(xdrI256.LoLo),
		}
		result.I256 = &i256
		return result, nil
	case xdr.ScValTypeScvBytes:
		xdrBytes := []byte(*v.Bytes)
		bytes := ScBytes(xdrBytes)
		result.Bytes = &bytes
		return result, nil
	case xdr.ScValTypeScvString:
		str := string(*v.Str)
		result.Str = &str
		return result, nil
	case xdr.ScValTypeScvSymbol:
		strSym := string(*v.Sym)
		sym := ScSymbol(strSym)
		result.Sym = &sym
		return result, nil
	case xdr.ScValTypeScvVec:
		xdrScVec := *v.Vec
		var ScVec []ScVal
		for _, xdrScVal := range *xdrScVec {
			scVal, err := ConvertScVal(xdrScVal)
			if err != nil {
				return result, err
			}
			ScVec = append(ScVec, scVal)
		}
		result.Vec = &ScVec
		return result, nil
	case xdr.ScValTypeScvMap:
		xdrScMap := *v.Map
		var scMapEntrys []ScMapEntry
		for _, xdrScMapEntry := range *xdrScMap {
			scMapEntry, err := ConvertScMapEntry(xdrScMapEntry)
			if err != nil {
				return result, err
			}

			scMapEntrys = append(scMapEntrys, scMapEntry)
		}
		scMap := ScMap(scMapEntrys)
		result.Map = &scMap
		return result, nil
	case xdr.ScValTypeScvAddress:
		xdrScAddress := *v.Address
		scAddress, err := ConvertScAddress(xdrScAddress)
		if err != nil {
			return result, err
		}
		result.Address = &scAddress
		return result, nil
	case xdr.ScValTypeScvContractInstance:
		xdrInstance := *v.Instance
		instance, err := ConvertScContractInstance(xdrInstance)
		if err != nil {
			return result, err
		}
		result.Instance = &instance
		return result, nil
	case xdr.ScValTypeScvLedgerKeyContractInstance:
		// void
		return result, nil
	case xdr.ScValTypeScvLedgerKeyNonce:
		xdrNonce := *v.NonceKey
		scNonce := ConvertScNonceKey(xdrNonce)
		result.NonceKey = &scNonce
		return result, nil
	}

	return result, errors.Errorf("error invalid ScVal type %v", v.Type)
}

func ConvertScMapEntry(m xdr.ScMapEntry) (ScMapEntry, error) {
	var result ScMapEntry

	key, err := ConvertScVal(m.Key)
	if err != nil {
		return result, err
	}

	val, err := ConvertScVal(m.Val)
	if err != nil {
		return result, err
	}

	result.Key = key
	result.Val = val

	return result, nil
}

func ConvertScContractInstance(i xdr.ScContractInstance) (ScContractInstance, error) {
	var result ScContractInstance
	executable, err := ConvertContractExecutable(i.Executable)
	if err != nil {

		return result, err
	}
	result.Executable = executable

	if i.Storage != nil {
		xdrStorage := *i.Storage
		var scMapEntrys []ScMapEntry
		for _, xdrScMapEntry := range xdrStorage {
			scMapEntry, err := ConvertScMapEntry(xdrScMapEntry)
			if err != nil {
				return result, err
			}
			scMapEntrys = append(scMapEntrys, scMapEntry)
		}
		storage := ScMap(scMapEntrys)
		result.Storage = &storage
	}

	return result, nil
}

func ConvertScNonceKey(k xdr.ScNonceKey) ScNonceKey {
	return ScNonceKey{
		Nonce: int64(k.Nonce),
	}
}

func ConvertExtensionPoint(p xdr.ExtensionPoint) ExtensionPoint {
	return ExtensionPoint{V: p.V}
}
