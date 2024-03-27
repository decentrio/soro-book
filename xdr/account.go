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
