package xdr

import (
	"time"

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

// TODO: testing
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

func ConvertClaimPredicates(inp []xdr.ClaimPredicate) ([]ClaimPredicate, error) {
	parts := make([]ClaimPredicate, len(inp))
	for i, pred := range inp {
		converted, err := ConvertClaimPredicate(pred)
		if err != nil {
			return parts, err
		}
		parts[i] = converted
	}
	return parts, nil
}

// TODO: testing
func ConvertClaimPredicate(cp xdr.ClaimPredicate) (ClaimPredicate, error) {
	var result ClaimPredicate

	switch cp.Type {
	case xdr.ClaimPredicateTypeClaimPredicateUnconditional:
		// void
		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateAnd:
		andPredicates, err := ConvertClaimPredicates(*cp.AndPredicates)
		if err != nil {
			return result, err
		}
		result.AndPredicates = &andPredicates

		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateOr:
		orPredicates, err := ConvertClaimPredicates(*cp.OrPredicates)
		if err != nil {
			return result, err
		}
		result.OrPredicates = &orPredicates

		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateNot:
		xdrNotPredicate, ok := cp.GetNotPredicate()
		if !ok {
			return result, errors.Errorf("invalid type ClaimPredicateTypeClaimPredicateNot")
		}

		notPredicate, err := ConvertClaimPredicate(*xdrNotPredicate)
		if err != nil {
			return result, err
		}
		result.NotPredicate = &notPredicate

		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateBeforeAbsoluteTime:
		absBeforeEpoch := int64(*cp.AbsBefore)
		absBefore := time.Unix(absBeforeEpoch, 0).UTC()

		result.AbsBefore = &absBefore
		result.AbsBeforeEpoch = &absBeforeEpoch

		return result, nil
	case xdr.ClaimPredicateTypeClaimPredicateBeforeRelativeTime:
		relBefore := int64(*cp.RelBefore)
		result.RelBefore = &relBefore

		return result, nil
	}

	return result, errors.Errorf("invalid ClaimPredicate type %v", cp.Type)
}

// TODO: testing
func ConvertClaimant(c xdr.Claimant) (Claimant, error) {
	var result Claimant

	switch c.Type {
	case xdr.ClaimantTypeClaimantTypeV0:
		xdrV0 := c.V0

		destination := PublicKey{
			Ed25519: xdrV0.Destination.Ed25519.String(),
		}

		predicate, err := ConvertClaimPredicate(c.V0.Predicate)
		if err != nil {
			return result, err
		}

		v0 := &ClaimantV0{
			Destination: destination,
			Predicate:   predicate,
		}
		result.V0 = v0

		return result, nil
	}

	return result, errors.Errorf("invalid claimant type %v", c.Type)
}

// TODO: testing
func ConvertClaimableBalanceId(id xdr.ClaimableBalanceId) (ClaimableBalanceId, error) {
	var result ClaimableBalanceId

	switch id.Type {
	case xdr.ClaimableBalanceIdTypeClaimableBalanceIdTypeV0:
		v0 := (*id.V0).HexString()
		result.V0 = &v0
		return result, nil
	}

	return result, errors.Errorf("invalid ClaimableBalanceId type %v", id.Type)
}

// TODO: testing
func ConvertPrice(p xdr.Price) Price {
	return Price{
		N: int32(p.N),
		D: int32(p.D),
	}
}
