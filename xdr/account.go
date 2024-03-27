package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

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
func ConvertEd25519(inp *xdr.Uint256) string {
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
