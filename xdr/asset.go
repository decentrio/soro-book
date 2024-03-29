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
func ConvertTrustLineAsset(a xdr.TrustLineAsset) (TrustLineAsset, error) {
	var result TrustLineAsset
	asset, err := ConvertAsset(a.ToAsset())
	if err != nil {
		return result, err
	}
	result.Asset = &asset

	xdrLpId := xdr.Hash(*a.LiquidityPoolId)
	lpId := PoolId(xdrLpId[:])
	result.LiquidityPoolId = &lpId

	return result, nil
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

func ConvertPathPaymentStrictReceiveResultSuccess(r xdr.PathPaymentStrictReceiveResultSuccess) (PathPaymentStrictReceiveResultSuccess, error) {
	var result PathPaymentStrictReceiveResultSuccess

	var offers []ClaimAtom
	for _, xdrOffer := range r.Offers {
		offer, err := ConvertClaimAtom(xdrOffer)
		if err != nil {
			return result, err
		}

		offers = append(offers, offer)
	}

	last, err := ConvertSimplePaymentResult(r.Last)
	if err != nil {
		return result, err
	}

	result.Offers = offers
	result.Last = last

	return result, nil
}

func ConvertPathPaymentStrictSendResultSuccess(r xdr.PathPaymentStrictSendResultSuccess) (PathPaymentStrictSendResultSuccess, error) {
	var result PathPaymentStrictSendResultSuccess

	var offers []ClaimAtom
	for _, xdrOffer := range r.Offers {
		offer, err := ConvertClaimAtom(xdrOffer)
		if err != nil {
			return result, err
		}

		offers = append(offers, offer)
	}

	last, err := ConvertSimplePaymentResult(r.Last)
	if err != nil {
		return result, err
	}

	result.Offers = offers
	result.Last = last

	return result, nil
}

func ConvertClaimAtom(c xdr.ClaimAtom) (ClaimAtom, error) {
	var result ClaimAtom

	switch c.Type {
	case xdr.ClaimAtomTypeClaimAtomTypeV0:
		v0, err := ConvertClaimOfferAtomV0(*c.V0)
		if err != nil {
			return result, err
		}

		result.V0 = &v0

		return result, nil
	case xdr.ClaimAtomTypeClaimAtomTypeOrderBook:
		orderBook, err := ConvertClaimOfferAtom(*c.OrderBook)
		if err != nil {
			return result, err
		}

		result.OrderBook = &orderBook

		return result, nil
	case xdr.ClaimAtomTypeClaimAtomTypeLiquidityPool:
		lp, err := ConvertClaimLiquidityAtom(*c.LiquidityPool)
		if err != nil {
			return result, err
		}

		result.LiquidityPool = &lp

		return result, nil
	}

	return result, errors.Errorf("invalid ConvertClaimAtom type %v", c.Type)
}

func ConvertClaimOfferAtomV0(c xdr.ClaimOfferAtomV0) (ClaimOfferAtomV0, error) {
	var result ClaimOfferAtomV0

	sellerEd25519 := c.SellerEd25519.String()

	assetSold, err := ConvertAsset(c.AssetSold)
	if err != nil {
		return result, err
	}

	assetBought, err := ConvertAsset(c.AssetBought)
	if err != nil {
		return result, err
	}

	result.SellerEd25519 = sellerEd25519
	result.OfferId = int64(c.OfferId)
	result.AssetSold = assetSold
	result.AmountSold = int64(c.AmountSold)
	result.AssetBought = assetBought
	result.AmountBought = int64(c.AmountBought)

	return result, nil
}

func ConvertClaimOfferAtom(c xdr.ClaimOfferAtom) (ClaimOfferAtom, error) {
	var result ClaimOfferAtom

	sellerId := PublicKey{
		Ed25519: c.SellerId.Ed25519.String(),
	}

	assetSold, err := ConvertAsset(c.AssetSold)
	if err != nil {
		return result, err
	}

	assetBought, err := ConvertAsset(c.AssetBought)
	if err != nil {
		return result, err
	}

	result.SellerId = sellerId
	result.OfferId = int64(c.OfferId)
	result.AssetSold = assetSold
	result.AmountSold = int64(c.AmountSold)
	result.AssetBought = assetBought
	result.AmountBought = int64(c.AmountBought)

	return result, nil
}

func ConvertClaimLiquidityAtom(c xdr.ClaimLiquidityAtom) (ClaimLiquidityAtom, error) {
	var result ClaimLiquidityAtom

	xdrPoolId := xdr.Hash(c.LiquidityPoolId)
	poolId := PoolId(xdrPoolId[:])

	assetSold, err := ConvertAsset(c.AssetSold)
	if err != nil {
		return result, err
	}

	assetBought, err := ConvertAsset(c.AssetBought)
	if err != nil {
		return result, err
	}

	result.LiquidityPoolId = poolId
	result.AssetSold = assetSold
	result.AmountSold = int64(c.AmountSold)
	result.AssetBought = assetBought
	result.AmountBought = int64(c.AmountBought)

	return result, nil
}

func ConvertSimplePaymentResult(r xdr.SimplePaymentResult) (SimplePaymentResult, error) {
	var result SimplePaymentResult
	destination := PublicKey{
		Ed25519: ConvertEd25519(r.Destination.Ed25519),
	}

	asset, err := ConvertAsset(r.Asset)
	if err != nil {
		return result, err
	}

	result.Destination = destination
	result.Asset = asset
	result.Amount = int64(r.Amount)

	return result, nil
}

func ConvertManageOfferSuccessResult(r xdr.ManageOfferSuccessResult) (ManageOfferSuccessResult, error) {
	var result ManageOfferSuccessResult

	var offersClaimed []ClaimAtom
	for _, xdrOffer := range r.OffersClaimed {
		offer, err := ConvertClaimAtom(xdrOffer)
		if err != nil {
			return result, err
		}

		offersClaimed = append(offersClaimed, offer)
	}

	offer, err := ConvertManageOfferSuccessResultOffer(r.Offer)
	if err != nil {
		return result, err
	}
	result.Offer = offer

	return result, nil
}

func ConvertManageOfferSuccessResultOffer(r xdr.ManageOfferSuccessResultOffer) (ManageOfferSuccessResultOffer, error) {
	var result ManageOfferSuccessResultOffer

	result.Effect = int32(r.Effect)

	offer, err := ConvertOfferEntry(*r.Offer)
	if err != nil {
		return result, err
	}

	result.Offer = &offer

	return result, nil
}

func ConvertOfferEntry(e xdr.OfferEntry) (OfferEntry, error) {
	var result OfferEntry

	sellerId := PublicKey{
		Ed25519: e.SellerId.Ed25519.String(),
	}

	selling, err := ConvertAsset(e.Selling)
	if err != nil {
		return result, err
	}

	buying, err := ConvertAsset(e.Buying)
	if err != nil {
		return result, err
	}

	price := ConvertPrice(e.Price)
	if err != nil {
		return result, err
	}

	result.SellerId = sellerId
	result.OfferId = int64(e.OfferId)
	result.Selling = selling
	result.Buying = buying
	result.Amount = int64(e.Amount)
	result.Price = price
	result.Flags = uint32(e.Flags)
	result.Ext = ConvertOfferEntryExt(e.Ext)

	return result, nil
}

func ConvertOfferEntryExt(e xdr.OfferEntryExt) OfferEntryExt {
	return OfferEntryExt{V: e.V}
}

func ConvertInflationPayout(i xdr.InflationPayout) InflationPayout {
	return InflationPayout{
		Destination: PublicKey{Ed25519: i.Destination.Ed25519.String()},
		Amount:      int64(i.Amount),
	}
}
