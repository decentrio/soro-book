package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

// TODO: testing
func ConvertAsset(as xdr.Asset) (Asset, error) {
	var result Asset
	switch as.Type {
	case xdr.AssetTypeAssetTypeCreditAlphanum4:
		result.AssetCode = as.AlphaNum4.AssetCode[:]

		issuer := PublicKey{
			Ed25519: ConvertEd25519(as.AlphaNum4.Issuer.ToMuxedAccount().Ed25519),
		}
		result.Issuer = issuer

		return result, nil
	case xdr.AssetTypeAssetTypeCreditAlphanum12:
		result.AssetCode = as.AlphaNum12.AssetCode[:]

		issuer := PublicKey{
			Ed25519: ConvertEd25519(as.AlphaNum12.Issuer.ToMuxedAccount().Ed25519),
		}
		result.Issuer = issuer

		return result, nil
	default:
		return result, errors.Errorf("unsupported asset type %v", as.Type)
	}

}

func ConvertPrice(p xdr.Price) Price {
	return Price{
		N: int32(p.N),
		D: int32(p.D),
	}
}
