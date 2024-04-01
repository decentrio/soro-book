package converter

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertLedgerEntryChange(c xdr.LedgerEntryChange) (LedgerEntryChange, error) {
	var result LedgerEntryChange

	switch c.Type {
	case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
		created, err := ConvertLedgerEntry(*c.Created)
		if err != nil {
			return result, err
		}

		result.Created = &created
		return result, nil
	case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
		updated, err := ConvertLedgerEntry(*c.Updated)
		if err != nil {
			return result, err
		}

		result.Updated = &updated
		return result, nil
	case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
		removed, err := ConvertLedgerKey(*c.Removed)
		if err != nil {
			return result, err
		}

		result.Removed = &removed
		return result, nil
	case xdr.LedgerEntryChangeTypeLedgerEntryState:
		state, err := ConvertLedgerEntry(*c.State)
		if err != nil {
			return result, err
		}

		result.State = &state
		return result, nil
	}
	return result, errors.Errorf("error invalid LedgerEntryChange type %v", c.Type)
}

func ConvertLedgerEntry(e xdr.LedgerEntry) (LedgerEntry, error) {
	var result LedgerEntry

	data, err := ConvertLedgerEntryData(e.Data)
	if err != nil {
		return result, err
	}

	ext := ConvertLedgerEntryExt(e.Ext)

	result.LastModifiedLedgerSeq = uint32(e.LastModifiedLedgerSeq)
	result.Data = data
	result.Ext = ext

	return result, nil
}

func ConvertLedgerEntryData(d xdr.LedgerEntryData) (LedgerEntryData, error) {
	var result LedgerEntryData
	switch d.Type {
	case xdr.LedgerEntryTypeAccount:
		account, err := ConvertAccountEntry(*d.Account)
		if err != nil {
			return result, err
		}
		result.Account = &account

		return result, nil
	case xdr.LedgerEntryTypeTrustline:
		trustLine, err := ConvertTrustLineEntry(*d.TrustLine)
		if err != nil {
			return result, err
		}
		result.TrustLine = &trustLine

		return result, nil
	case xdr.LedgerEntryTypeOffer:
		offer, err := ConvertOfferEntry(*d.Offer)
		if err != nil {
			return result, err
		}
		result.Offer = &offer

		return result, nil
	case xdr.LedgerEntryTypeData:
		data := ConvertDataEntry(*d.Data)
		result.Data = &data
		return result, nil
	case xdr.LedgerEntryTypeClaimableBalance:
		balance, err := ConvertConvertClaimableBalanceEntry(*d.ClaimableBalance)
		if err != nil {
			return result, err
		}
		result.ClaimableBalance = &balance

		return result, nil
	case xdr.LedgerEntryTypeLiquidityPool:
		lp, err := ConvertLiquidityPoolEntry(*d.LiquidityPool)
		if err != nil {
			return result, err
		}
		result.LiquidityPool = &lp

		return result, nil
	case xdr.LedgerEntryTypeContractData:
		contractData, err := ConvertContractDataEntry(*d.ContractData)
		if err != nil {
			return result, err
		}
		result.ContractData = &contractData

		return result, nil
	case xdr.LedgerEntryTypeContractCode:
		contractCode := ConvertContractCodeEntry(*d.ContractCode)
		result.ContractCode = &contractCode

		return result, nil
	case xdr.LedgerEntryTypeConfigSetting:
		cfgSettings, err := ConvertConfigSettingEntry(*d.ConfigSetting)
		if err != nil {
			return result, err
		}
		result.ConfigSetting = &cfgSettings

		return result, nil
	case xdr.LedgerEntryTypeTtl:
		ttl := ConvertTtlEntry(*d.Ttl)
		result.Ttl = &ttl

		return result, nil
	}

	return result, errors.Errorf("error invalid LedgerEntryData type %v", d.Type)
}

func ConvertLedgerEntryExt(e xdr.LedgerEntryExt) LedgerEntryExt {
	var v1 LedgerEntryExtensionV1
	if e.V1 != nil {
		v1 = ConvertLedgerEntryExtensionV1(*e.V1)
	}

	return LedgerEntryExt{
		V:  e.V,
		V1: &v1,
	}
}

func ConvertLedgerEntryExtensionV1(e xdr.LedgerEntryExtensionV1) LedgerEntryExtensionV1 {
	var sponsoringId PublicKey
	if e.SponsoringId != nil {
		sponsoringId = PublicKey{
			Ed25519: e.SponsoringId.Ed25519.String(),
		}
	}

	return LedgerEntryExtensionV1{
		SponsoringId: sponsoringId,
		Ext:          ConvertLedgerEntryExtensionV1Ext(e.Ext),
	}
}

func ConvertLedgerEntryExtensionV1Ext(e xdr.LedgerEntryExtensionV1Ext) LedgerEntryExtensionV1Ext {
	return LedgerEntryExtensionV1Ext{V: e.V}
}

func ConvertLedgerKeyAccount(k xdr.LedgerKeyAccount) LedgerKeyAccount {
	accountId := PublicKey{
		Ed25519: ConvertEd25519(k.AccountId.Ed25519),
	}

	return LedgerKeyAccount{
		AccountId: accountId,
	}
}

func ConvertLedgerKeyTrustLine(k xdr.LedgerKeyTrustLine) (result LedgerKeyTrustLine, err error) {
	accountID := PublicKey{
		Ed25519: ConvertEd25519(k.AccountId.Ed25519),
	}

	asset, err := ConvertTrustLineAsset(k.Asset)
	if err != nil {
		return result, err
	}

	result.AccountId = accountID
	result.Asset = asset

	return result, nil
}

func ConvertLedgerKeyOffer(k xdr.LedgerKeyOffer) (result LedgerKeyOffer, err error) {
	seller := PublicKey{
		Ed25519: ConvertEd25519(k.SellerId.Ed25519),
	}

	result.SellerId = seller
	result.OfferId = int64(k.OfferId)

	return result, err
}

func ConvertLedgerKeyData(k xdr.LedgerKeyData) (LedgerKeyData, error) {
	var result LedgerKeyData

	accountID := PublicKey{
		Ed25519: ConvertEd25519(k.AccountId.Ed25519),
	}

	result.AccountId = accountID
	result.DataName = string(k.DataName)

	return result, nil
}

func ConvertLedgerKeyClaimableBalance(k xdr.LedgerKeyClaimableBalance) (LedgerKeyClaimableBalance, error) {
	var result LedgerKeyClaimableBalance

	id, err := ConvertClaimableBalanceId(k.BalanceId)
	if err != nil {
		return result, err
	}

	result.BalanceId = id

	return result, nil
}

func ConvertLedgerKeyLiquidityPool(k xdr.LedgerKeyLiquidityPool) (LedgerKeyLiquidityPool, error) {
	var result LedgerKeyLiquidityPool

	xdrHashPoolId := xdr.Hash(k.LiquidityPoolId)
	lpId := PoolId(xdrHashPoolId[:])
	result.LiquidityPoolId = lpId

	return result, nil
}

func ConvertLedgerKeyContractData(k xdr.LedgerKeyContractData) (LedgerKeyContractData, error) {
	var result LedgerKeyContractData

	contract, err := ConvertScAddress(k.Contract)
	if err != nil {
		return result, err
	}
	result.Contract = contract

	key, err := ConvertScVal(k.Key)
	if err != nil {
		return result, err
	}
	result.Key = key

	result.Durability = int32(k.Durability)

	return result, nil
}

func ConvertLedgerKeyContractCode(k xdr.LedgerKeyContractCode) (LedgerKeyContractCode, error) {
	var result LedgerKeyContractCode
	result.Hash = k.Hash.HexString()

	return result, nil
}

func ConvertLedgerKeyConfigSetting(k xdr.LedgerKeyConfigSetting) (LedgerKeyConfigSetting, error) {
	var result LedgerKeyConfigSetting
	result.ConfigSettingId = int32(k.ConfigSettingId)

	return result, nil
}

func ConvertLedgerKeyTtl(k xdr.LedgerKeyTtl) (LedgerKeyTtl, error) {
	var result LedgerKeyTtl
	result.KeyHash = k.KeyHash.HexString()

	return result, nil
}

func ConvertLedgerFootprint(f xdr.LedgerFootprint) (LedgerFootprint, error) {
	var result LedgerFootprint

	var readOnlys []LedgerKey
	for _, ledgerKey := range f.ReadOnly {
		readOnly, err := ConvertLedgerKey(ledgerKey)
		if err != nil {
			return result, err
		}

		readOnlys = append(readOnlys, readOnly)
	}

	var readWrites []LedgerKey
	for _, ledgerKey := range f.ReadWrite {
		readWrite, err := ConvertLedgerKey(ledgerKey)
		if err != nil {
			return result, err
		}

		readWrites = append(readWrites, readWrite)
	}

	result.ReadOnly = readOnlys
	result.ReadWrite = readWrites

	return result, nil
}

// TODO: testing
func ConvertLedgerKey(k xdr.LedgerKey) (LedgerKey, error) {
	var result LedgerKey
	switch k.Type {
	case xdr.LedgerEntryTypeAccount:
		account := ConvertLedgerKeyAccount(*k.Account)
		result.Account = &account
		return result, nil
	case xdr.LedgerEntryTypeTrustline:
		trustLine, err := ConvertLedgerKeyTrustLine(*k.TrustLine)
		if err != nil {
			return result, err
		}
		result.TrustLine = &trustLine
		return result, nil
	case xdr.LedgerEntryTypeOffer:
		offer, err := ConvertLedgerKeyOffer(*k.Offer)
		if err != nil {
			return result, err
		}
		result.Offer = &offer
		return result, nil
	case xdr.LedgerEntryTypeData:
		data, err := ConvertLedgerKeyData(*k.Data)
		if err != nil {
			return result, err
		}
		result.Data = &data
		return result, nil
	case xdr.LedgerEntryTypeClaimableBalance:
		claimableBalance, err := ConvertLedgerKeyClaimableBalance(*k.ClaimableBalance)
		if err != nil {
			return result, err
		}
		result.ClaimableBalance = &claimableBalance
		return result, nil
	case xdr.LedgerEntryTypeLiquidityPool:
		liquidityPool, err := ConvertLedgerKeyLiquidityPool(*k.LiquidityPool)
		if err != nil {
			return result, err
		}
		result.LiquidityPool = &liquidityPool
		return result, nil
	case xdr.LedgerEntryTypeContractData:
		contractData, err := ConvertLedgerKeyContractData(*k.ContractData)
		if err != nil {
			return result, err
		}
		result.ContractData = &contractData
		return result, nil
	case xdr.LedgerEntryTypeContractCode:
		contractCode, err := ConvertLedgerKeyContractCode(*k.ContractCode)
		if err != nil {
			return result, err
		}
		result.ContractCode = &contractCode
		return result, nil
	case xdr.LedgerEntryTypeConfigSetting:
		cfgSetting, err := ConvertLedgerKeyConfigSetting(*k.ConfigSetting)
		if err != nil {
			return result, err
		}
		result.ConfigSetting = &cfgSetting
		return result, nil
	case xdr.LedgerEntryTypeTtl:
		ttl, err := ConvertLedgerKeyTtl(*k.Ttl)
		if err != nil {
			return result, err
		}
		result.Ttl = &ttl
		return result, nil
	}

	return result, errors.Errorf("error invalid LedgerKey type %v", k.Type)
}

func ConvertLedgerBounds(b xdr.LedgerBounds) LedgerBounds {
	return LedgerBounds{
		MinLedger: uint32(b.MinLedger),
		MaxLedger: uint32(b.MaxLedger),
	}
}
