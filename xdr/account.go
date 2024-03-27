package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

// TODO: testing
func ConvertSigner(s xdr.Signer) (Signer, error) {
	var result Signer
	signerKey, err := ConvertSignerKey(s.Key)
	if err != nil {
		return result, err
	}
	result.Key = signerKey
	result.Weight = uint32(s.Weight)

	return result, nil
}

// TODO: testing
func ConvertSignerKey(k xdr.SignerKey) (SignerKey, error) {
	var result SignerKey
	switch k.Type {
	case xdr.SignerKeyTypeSignerKeyTypeEd25519:
		ed25519 := ConvertEd25519(k.Ed25519)
		result.Ed25519 = &ed25519
		return result, nil
	case xdr.SignerKeyTypeSignerKeyTypePreAuthTx:
		preAuthTx := ConvertPreAuthTx(k.PreAuthTx)
		result.PreAuthTx = &preAuthTx
		return result, nil
	case xdr.SignerKeyTypeSignerKeyTypeHashX:
		hashX := ConvertHashX(k.HashX)
		result.HashX = &hashX
		return result, nil
	case xdr.SignerKeyTypeSignerKeyTypeEd25519SignedPayload:
		signedPayload := ConvertSignerKeyEd25519SignedPayload(k.Ed25519SignedPayload)
		result.Ed25519SignedPayload = &signedPayload
		return result, nil
	}

	return result, errors.Errorf("error invalid signer key type %v", k.Type)
}

// TODO: testing
func ConvertMuxedAccount(ma xdr.MuxedAccount) (MuxedAccount, error) {
	var result MuxedAccount
	switch ma.Type {
	case xdr.CryptoKeyTypeKeyTypeEd25519:
		key := ConvertEd25519(ma.Ed25519)
		result.Ed25519 = &key

		return result, nil
	case xdr.CryptoKeyTypeKeyTypeMuxedEd25519:
		mam := ConvertMuxedAccountMed25519(ma.Med25519)
		result.Med25519 = &mam

		return result, nil
	}

	return MuxedAccount{}, errors.Errorf("error invalid muxed account type %v", ma.Type)
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
}

func ConvertLedgerKeyContractCode(k xdr.LedgerKeyContractCode) (LedgerKeyContractCode, error) {
	var result LedgerKeyContractCode
}

func ConvertLedgerKeyConfigSetting(k xdr.LedgerKeyConfigSetting) (LedgerKeyConfigSetting, error) {
	var result LedgerKeyConfigSetting
}

func ConvertLedgerKeyTtl(k xdr.LedgerKeyTtl) (LedgerKeyTtl, error) {
	var result LedgerKeyTtl
}

// TODO: testing
func ConvertLegerKey(k xdr.LedgerKey) (LedgerKey, error) {
	var result LedgerKey
}

// TODO :testing
func ConvertRevokeSponsorshipOpSigner(s xdr.RevokeSponsorshipOpSigner) (RevokeSponsorshipOpSigner, error) {
	var result RevokeSponsorshipOpSigner
}

// TODO: testing
func ConvertSignerKeyEd25519SignedPayload(inp *xdr.SignerKeyEd25519SignedPayload) SignerKeyEd25519SignedPayload {
	result := SignerKeyEd25519SignedPayload{
		Ed25519: inp.Ed25519.String(),
		Payload: inp.Payload,
	}
	return result
}

// TODO: testing
func ConvertEd25519(inp *xdr.Uint256) string {
	result := inp.String()
	return result
}

// TODO: testing
func ConvertPreAuthTx(inp *xdr.Uint256) string {
	result := inp.String()
	return result
}

// TODO: testing
func ConvertHashX(inp *xdr.Uint256) string {
	result := inp.String()
	return result
}

// TODO: testing
func ConvertMuxedAccountMed25519(inp *xdr.MuxedAccountMed25519) MuxedAccountMed25519 {
	key := ConvertEd25519(&inp.Ed25519)
	return MuxedAccountMed25519{
		Id:      uint64(inp.Id),
		Ed25519: &key,
	}
}
