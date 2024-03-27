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

func ConvertLiquidityPoolConstantProductParameters(
	lpcpp xdr.LiquidityPoolConstantProductParameters,
) (LiquidityPoolConstantProductParameters, error) {
	var result LiquidityPoolConstantProductParameters

	assetA, err := ConvertAsset(lpcpp.AssetA)
	if err != nil {
		return result, err
	}

	assetB, err := ConvertAsset(lpcpp.AssetB)
	if err != nil {
		return result, err
	}

	result.AssetA = assetA
	result.AssetB = assetB
	result.Fee = int32(lpcpp.Fee)

	return result, nil
}

// TODO: testing
func ConvertLiquidityPoolParameters(lpp xdr.LiquidityPoolParameters) (LiquidityPoolParameters, error) {
	var result LiquidityPoolParameters

	switch lpp.Type {
	case xdr.LiquidityPoolTypeLiquidityPoolConstantProduct:
		lpcpp, err := ConvertLiquidityPoolConstantProductParameters(*lpp.ConstantProduct)
		if err != nil {
			return result, err
		}

		result.ConstantProduct = &lpcpp

		return result, nil
	}

	return result, errors.Errorf("invalid liquidity pool parameters type %v", lpp.Type)
}

// TODO: testing
func ConvertChangeTrustAsset(ta xdr.ChangeTrustAsset) (ChangeTrustAsset, error) {
	var result ChangeTrustAsset

	asset, err := ConvertAsset(ta.ToAsset())
	if err != nil {
		return result, err
	}
	result.Asset = &asset

	liquidityPool, err := ConvertLiquidityPoolParameters(*ta.LiquidityPool)
	if err != nil {
		return result, err
	}
	result.LiquidityPool = &liquidityPool

	return result, nil
}

// TODO: testing
func ConvertPrice(p xdr.Price) Price {
	return Price{
		N: int32(p.N),
		D: int32(p.D),
	}
}
