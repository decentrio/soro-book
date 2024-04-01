package converter

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertAccountEntry(e xdr.AccountEntry) (AccountEntry, error) {
	var result AccountEntry

	accountId := PublicKey{
		Ed25519: e.AccountId.Ed25519.String(),
	}

	var inflationDest PublicKey
	if e.InflationDest != nil {
		inflationDest = PublicKey{
			Ed25519: (*e.InflationDest).Ed25519.String(),
		}
	}

	var signers []Signer
	for _, xdrSigner := range e.Signers {
		signer, err := ConvertSigner(xdrSigner)
		if err != nil {
			return result, err
		}

		signers = append(signers, signer)
	}

	ext := ConvertAccountEntryExt(e.Ext)

	result.AccountId = accountId
	result.Balance = int64(e.Balance)
	result.SeqNum = int64(e.SeqNum)
	result.NumSubEntries = uint32(e.NumSubEntries)
	result.InflationDest = &inflationDest
	result.Flags = uint32(e.Flags)
	result.HomeDomain = string(e.HomeDomain)
	result.Thresholds = e.Thresholds[:]
	result.Signers = signers
	result.Ext = ext

	return result, nil
}

func ConvertAccountEntryExt(e xdr.AccountEntryExt) AccountEntryExt {
	var v1 AccountEntryExtensionV1
	if e.V1 != nil {
		v1 = ConvertAccountEntryExtensionV1(*e.V1)
	}

	return AccountEntryExt{
		V:  e.V,
		V1: &v1,
	}
}

func ConvertAccountEntryExtensionV1(e xdr.AccountEntryExtensionV1) AccountEntryExtensionV1 {
	return AccountEntryExtensionV1{
		Liabilities: ConvertLiabilities(e.Liabilities),
		Ext:         ConvertAccountEntryExtensionV1Ext(e.Ext),
	}
}

func ConvertLiabilities(l xdr.Liabilities) Liabilities {
	return Liabilities{
		Buying:  int64(l.Buying),
		Selling: int64(l.Selling),
	}
}

func ConvertAccountEntryExtensionV1Ext(e xdr.AccountEntryExtensionV1Ext) AccountEntryExtensionV1Ext {
	var v2 AccountEntryExtensionV2
	if e.V2 != nil {
		v2 = ConvertAccountEntryExtensionV2(*e.V2)
	}

	return AccountEntryExtensionV1Ext{
		V:  e.V,
		V2: &v2,
	}
}

func ConvertAccountEntryExtensionV2(e xdr.AccountEntryExtensionV2) AccountEntryExtensionV2 {
	var signerSponsoringIDs []PublicKey
	if e.SignerSponsoringIDs != nil {
		for _, xdrSigner := range e.SignerSponsoringIDs {
			if xdrSigner != nil {
				signer := PublicKey{
					Ed25519: (*xdrSigner.Ed25519).String(),
				}

				signerSponsoringIDs = append(signerSponsoringIDs, signer)
			}

		}
	}

	ext := ConvertAccountEntryExtensionV2Ext(e.Ext)

	return AccountEntryExtensionV2{
		NumSponsored:        uint32(e.NumSponsored),
		NumSponsoring:       uint32(e.NumSponsoring),
		SignerSponsoringIDs: signerSponsoringIDs,
		Ext:                 ext,
	}
}

func ConvertAccountEntryExtensionV2Ext(e xdr.AccountEntryExtensionV2Ext) AccountEntryExtensionV2Ext {
	var v3 AccountEntryExtensionV3
	if e.V3 != nil {
		v3 = ConvertAccountEntryExtensionV3(*e.V3)
	}

	return AccountEntryExtensionV2Ext{
		V:  e.V,
		V3: &v3,
	}
}

func ConvertAccountEntryExtensionV3(e xdr.AccountEntryExtensionV3) AccountEntryExtensionV3 {
	return AccountEntryExtensionV3{
		Ext:       ConvertExtensionPoint(e.Ext),
		SeqLedger: uint32(e.SeqLedger),
		SeqTime:   uint64(e.SeqTime),
	}
}

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

// TODO :testing
func ConvertRevokeSponsorshipOpSigner(s xdr.RevokeSponsorshipOpSigner) (RevokeSponsorshipOpSigner, error) {
	var result RevokeSponsorshipOpSigner

	accountId := PublicKey{
		Ed25519: s.AccountId.Ed25519.String(),
	}

	signerKey, err := ConvertSignerKey(s.SignerKey)
	if err != nil {
		return result, err
	}

	result.AccountId = accountId
	result.SignerKey = signerKey

	return result, nil
}

func ConvertDecoratedSignature(s xdr.DecoratedSignature) DecoratedSignature {
	return DecoratedSignature{
		Hint:      s.Hint[:],
		Signature: s.Signature,
	}
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
