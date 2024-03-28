package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

// TODO: testing
func ConvertOperation(op xdr.Operation) (Operation, error) {
	var result Operation
	sourceAccount, err := ConvertMuxedAccount(*op.SourceAccount)
	if err != nil {
		return result, err
	}
	result.SourceAccount = &sourceAccount

	body, err := ConvertOperationBody(op.Body)
	if err != nil {
		return result, err
	}
	result.Body = body

	return result, nil
}

// TODO: testing
func ConvertOperationBody(bd xdr.OperationBody) (OperationBody, error) {
	var result OperationBody

	switch bd.Type {
	case xdr.OperationTypeCreateAccount:
		xdrDestination := bd.CreateAccountOp.Destination
		destination := PublicKey{
			Ed25519: xdrDestination.Ed25519.String(),
		}
		createAccountOp := &CreateAccountOp{
			Destination:     destination,
			StartingBalance: int64(bd.CreateAccountOp.StartingBalance),
		}
		result.CreateAccountOp = createAccountOp

		return result, nil
	case xdr.OperationTypePayment:
		xdrPaymentOp := bd.PaymentOp

		destination, err := ConvertMuxedAccount(xdrPaymentOp.Destination)
		if err != nil {
			return result, err
		}

		asset, err := ConvertAsset(xdrPaymentOp.Asset)
		if err != nil {
			return result, err
		}

		paymentOp := &PaymentOp{
			Destination: destination,
			Asset:       asset,
			Amount:      int64(xdrPaymentOp.Amount),
		}
		result.PaymentOp = paymentOp

		return result, nil
	case xdr.OperationTypePathPaymentStrictReceive:
		xdrPathPaymentStrictReceiveOp := bd.PathPaymentStrictReceiveOp

		sendAsset, err := ConvertAsset(xdrPathPaymentStrictReceiveOp.SendAsset)
		if err != nil {
			return result, err
		}

		destination, err := ConvertMuxedAccount(xdrPathPaymentStrictReceiveOp.Destination)
		if err != nil {
			return result, err
		}

		destAsset, err := ConvertAsset(xdrPathPaymentStrictReceiveOp.DestAsset)
		if err != nil {
			return result, err
		}

		var paths []Asset
		for _, xdrPath := range xdrPathPaymentStrictReceiveOp.Path {
			path, err := ConvertAsset(xdrPath)
			if err != nil {
				return result, err
			}

			paths = append(paths, path)
		}

		pathPaymentStrictReceiveOp := &PathPaymentStrictReceiveOp{
			SendAsset:   sendAsset,
			SendMax:     int64(xdrPathPaymentStrictReceiveOp.SendMax),
			Destination: destination,
			DestAsset:   destAsset,
			DestAmount:  int64(xdrPathPaymentStrictReceiveOp.DestAmount),
			Path:        paths,
		}
		result.PathPaymentStrictReceiveOp = pathPaymentStrictReceiveOp

		return result, nil
	case xdr.OperationTypeManageSellOffer:
		xdrManageSellOffer := bd.ManageBuyOfferOp

		selling, err := ConvertAsset(xdrManageSellOffer.Selling)
		if err != nil {
			return result, err
		}

		buying, err := ConvertAsset(xdrManageSellOffer.Buying)
		if err != nil {
			return result, err
		}

		price := ConvertPrice(xdrManageSellOffer.Price)

		managerSellOfferOp := &ManageSellOfferOp{
			Selling:   selling,
			Buying:    buying,
			BuyAmount: int64(xdrManageSellOffer.BuyAmount),
			Price:     price,
			OfferId:   int64(xdrManageSellOffer.OfferId),
		}
		result.ManageSellOfferOp = managerSellOfferOp

		return result, nil
	case xdr.OperationTypeCreatePassiveSellOffer:
		xdrCreatePassiveSellOffer := bd.CreatePassiveSellOfferOp

		selling, err := ConvertAsset(xdrCreatePassiveSellOffer.Selling)
		if err != nil {
			return result, err
		}

		buying, err := ConvertAsset(xdrCreatePassiveSellOffer.Buying)
		if err != nil {
			return result, err
		}

		price := ConvertPrice(xdrCreatePassiveSellOffer.Price)
		createPassiveSellOffer := &CreatePassiveSellOfferOp{
			Selling: selling,
			Buying:  buying,
			Amount:  int64(xdrCreatePassiveSellOffer.Amount),
			Price:   price,
		}
		result.CreatePassiveSellOfferOp = createPassiveSellOffer

		return result, nil
	case xdr.OperationTypeSetOptions:
		xdrSetOptions := bd.SetOptionsOp
		inflationDest := PublicKey{
			Ed25519: ConvertEd25519(xdrSetOptions.InflationDest.Ed25519),
		}

		clearFlags := uint32(*xdrSetOptions.ClearFlags)
		setFlags := uint32(*xdrSetOptions.SetFlags)
		masterWeight := uint32(*xdrSetOptions.MasterWeight)
		lowThreshold := uint32(*xdrSetOptions.LowThreshold)
		medThreshold := uint32(*xdrSetOptions.MedThreshold)
		highThreshold := uint32(*xdrSetOptions.HighThreshold)
		homeDomain := string(*xdrSetOptions.HomeDomain)

		signer, err := ConvertSigner(*xdrSetOptions.Signer)
		if err != nil {
			return result, err
		}

		setOptions := &SetOptionsOp{
			InflationDest: &inflationDest,
			ClearFlags:    &clearFlags,
			SetFlags:      &setFlags,
			MasterWeight:  &masterWeight,
			LowThreshold:  &lowThreshold,
			MedThreshold:  &medThreshold,
			HighThreshold: &highThreshold,
			HomeDomain:    &homeDomain,
			Signer:        &signer,
		}
		result.SetOptionsOp = setOptions

		return result, nil
	case xdr.OperationTypeChangeTrust:
		xdrChangeTrust := bd.ChangeTrustOp

		line, err := ConvertChangeTrustAsset(xdrChangeTrust.Line)
		if err != nil {
			return result, err
		}

		changeTrust := &ChangeTrustOp{
			Line:  line,
			Limit: int64(xdrChangeTrust.Limit),
		}
		result.ChangeTrustOp = changeTrust

		return result, nil
	case xdr.OperationTypeAllowTrust:
		xdrAllowTrust := bd.AllowTrustOp

		trustor := PublicKey{
			Ed25519: ConvertEd25519(xdrAllowTrust.Trustor.Ed25519),
		}

		var assetCode []byte
		switch xdrAllowTrust.Asset.Type {
		case xdr.AssetTypeAssetTypeCreditAlphanum4:
			assetCode = xdrAllowTrust.Asset.AssetCode4[:]
		case xdr.AssetTypeAssetTypeCreditAlphanum12:
			assetCode = xdrAllowTrust.Asset.AssetCode12[:]
		default:
			return result, errors.Errorf("OperationTypeAllowTrust invalid asset code %v", xdrAllowTrust.Asset.Type)
		}

		allowTrust := &AllowTrustOp{
			Trustor:   trustor,
			AssetCode: assetCode,
			Authorize: uint32(xdrAllowTrust.Authorize),
		}
		result.AllowTrustOp = allowTrust

		return result, nil
	case xdr.OperationTypeAccountMerge:
		xdrDestination := bd.Destination
		destination, err := ConvertMuxedAccount(*xdrDestination)
		if err != nil {
			return result, err
		}
		result.Destination = &destination

		return result, nil
	case xdr.OperationTypeInflation:
		// void
		return result, nil
	case xdr.OperationTypeManageData:
		xdrManageDataOp := bd.ManageDataOp

		mangeData := &ManageDataOp{
			DataName:  string(xdrManageDataOp.DataName),
			DataValue: *xdrManageDataOp.DataValue,
		}
		result.ManageDataOp = mangeData

		return result, nil
	case xdr.OperationTypeBumpSequence:
		xdrBumpSequenceOp := bd.BumpSequenceOp

		bumpSequenceOp := &BumpSequenceOp{
			BumpTo: int64(xdrBumpSequenceOp.BumpTo),
		}
		result.BumpSequenceOp = bumpSequenceOp

		return result, nil
	case xdr.OperationTypeManageBuyOffer:
		xdrManageBuyOfferOp := bd.ManageBuyOfferOp

		selling, err := ConvertAsset(xdrManageBuyOfferOp.Selling)
		if err != nil {
			return result, err
		}

		buying, err := ConvertAsset(xdrManageBuyOfferOp.Buying)
		if err != nil {
			return result, err
		}

		price := ConvertPrice(xdrManageBuyOfferOp.Price)

		manageBuyOfferOp := &ManageBuyOfferOp{
			Selling:   selling,
			Buying:    buying,
			BuyAmount: int64(xdrManageBuyOfferOp.BuyAmount),
			Price:     price,
			OfferId:   int64(xdrManageBuyOfferOp.OfferId),
		}
		result.ManageBuyOfferOp = manageBuyOfferOp

		return result, nil
	case xdr.OperationTypePathPaymentStrictSend:
		xdrPathPaymentStrictSendOp := bd.PathPaymentStrictSendOp

		sendAsset, err := ConvertAsset(xdrPathPaymentStrictSendOp.SendAsset)
		if err != nil {
			return result, err
		}

		destAsset, err := ConvertAsset(xdrPathPaymentStrictSendOp.DestAsset)
		if err != nil {
			return result, err
		}

		var paths []Asset
		for _, xdrPath := range xdrPathPaymentStrictSendOp.Path {
			path, err := ConvertAsset(xdrPath)
			if err != nil {
				return result, err
			}

			paths = append(paths, path)
		}

		destination, err := ConvertMuxedAccount(xdrPathPaymentStrictSendOp.Destination)
		if err != nil {
			return result, err
		}

		pathPaymentStrictSendOp := &PathPaymentStrictSendOp{
			SendAsset:   sendAsset,
			SendAmount:  int64(xdrPathPaymentStrictSendOp.SendAmount),
			Destination: destination,
			DestAsset:   destAsset,
			DestMin:     int64(xdrPathPaymentStrictSendOp.DestMin),
			Path:        paths,
		}
		result.PathPaymentStrictSendOp = pathPaymentStrictSendOp

		return result, nil
	case xdr.OperationTypeCreateClaimableBalance:
		xdrCreateClaimableBalanceOp := bd.CreateClaimableBalanceOp

		asset, err := ConvertAsset(xdrCreateClaimableBalanceOp.Asset)
		if err != nil {
			return result, nil
		}

		var claimaints []Claimant
		for _, xdrClaimant := range xdrCreateClaimableBalanceOp.Claimants {
			claimant, err := ConvertClaimant(xdrClaimant)
			if err != nil {
				return result, nil
			}

			claimaints = append(claimaints, claimant)
		}

		createClaimableBalanceOp := &CreateClaimableBalanceOp{
			Asset:     asset,
			Amount:    int64(xdrCreateClaimableBalanceOp.Amount),
			Claimants: claimaints,
		}
		result.CreateClaimableBalanceOp = createClaimableBalanceOp

		return result, nil
	case xdr.OperationTypeClaimClaimableBalance:
		xdrClaimClaimableBalanceOp := bd.ClaimClaimableBalanceOp

		balanceId, err := ConvertClaimableBalanceId(xdrClaimClaimableBalanceOp.BalanceId)
		if err != nil {
			return result, nil
		}

		claimClaimableBalanceOp := &ClaimClaimableBalanceOp{
			BalanceId: balanceId,
		}
		result.ClaimClaimableBalanceOp = claimClaimableBalanceOp

		return result, nil
	case xdr.OperationTypeBeginSponsoringFutureReserves:
		xdrBeginSponsoringFutureReservesOp := bd.BeginSponsoringFutureReservesOp

		sponsoredId := PublicKey{
			Ed25519: ConvertEd25519(xdrBeginSponsoringFutureReservesOp.SponsoredId.Ed25519),
		}

		beginSponsoringFutureReservesOp := &BeginSponsoringFutureReservesOp{
			SponsoredId: sponsoredId,
		}
		result.BeginSponsoringFutureReservesOp = beginSponsoringFutureReservesOp

		return result, nil
	case xdr.OperationTypeEndSponsoringFutureReserves:
		// void
		return result, nil
	case xdr.OperationTypeRevokeSponsorship:
		xdrRevokeSponsorshipOp := bd.RevokeSponsorshipOp

		ledgerKey, err := ConvertLedgerKey(*xdrRevokeSponsorshipOp.LedgerKey)
		if err != nil {
			return result, nil
		}

		signer, err := ConvertRevokeSponsorshipOpSigner(*xdrRevokeSponsorshipOp.Signer)
		if err != nil {
			return result, nil
		}

		revokeSponsorshipOp := &RevokeSponsorshipOp{
			LedgerKey: &ledgerKey,
			Signer:    &signer,
		}
		result.RevokeSponsorshipOp = revokeSponsorshipOp

		return result, nil
	case xdr.OperationTypeClawback:
		xdrClawbackOp := bd.ClawbackOp

		asset, err := ConvertAsset(xdrClawbackOp.Asset)
		if err != nil {
			return result, nil
		}

		from, err := ConvertMuxedAccount(xdrClawbackOp.From)
		if err != nil {
			return result, nil
		}

		clawbackOp := &ClawbackOp{
			Asset:  asset,
			From:   from,
			Amount: int64(xdrClawbackOp.Amount),
		}
		result.ClawbackOp = clawbackOp

		return result, nil
	case xdr.OperationTypeClawbackClaimableBalance:
		xdrClawbackClaimableBalanceOp := bd.ClawbackClaimableBalanceOp

		balanceId, err := ConvertClaimableBalanceId(xdrClawbackClaimableBalanceOp.BalanceId)
		if err != nil {
			return result, nil
		}

		clawbackClaimableBalanceOp := &ClawbackClaimableBalanceOp{
			BalanceId: balanceId,
		}
		result.ClawbackClaimableBalanceOp = clawbackClaimableBalanceOp

		return result, nil
	case xdr.OperationTypeSetTrustLineFlags:
		xdrSetTrustLineFlagsOp := bd.SetTrustLineFlagsOp

		trustor := PublicKey{
			Ed25519: ConvertEd25519(xdrSetTrustLineFlagsOp.Trustor.Ed25519),
		}

		asset, err := ConvertAsset(xdrSetTrustLineFlagsOp.Asset)
		if err != nil {
			return result, nil
		}

		setTrustLineFlagsOp := &SetTrustLineFlagsOp{
			Trustor:    trustor,
			Asset:      asset,
			ClearFlags: uint32(xdrSetTrustLineFlagsOp.ClearFlags),
			SetFlags:   uint32(xdrSetTrustLineFlagsOp.SetFlags),
		}
		result.SetTrustLineFlagsOp = setTrustLineFlagsOp

		return result, nil
	case xdr.OperationTypeLiquidityPoolDeposit:
		xdrLiquidityPoolDepositOp := bd.LiquidityPoolDepositOp

		xdrHashLpId := xdr.Hash(xdrLiquidityPoolDepositOp.LiquidityPoolId)
		LpId := PoolId(xdrHashLpId[:])

		minPrice := ConvertPrice(xdrLiquidityPoolDepositOp.MinPrice)
		maxPrice := ConvertPrice(xdrLiquidityPoolDepositOp.MaxPrice)

		liquidityPoolDepositOp := &LiquidityPoolDepositOp{
			LiquidityPoolId: LpId,
			MaxAmountA:      int64(xdrLiquidityPoolDepositOp.MaxAmountA),
			MaxAmountB:      int64(xdrLiquidityPoolDepositOp.MaxAmountB),
			MinPrice:        minPrice,
			MaxPrice:        maxPrice,
		}
		result.LiquidityPoolDepositOp = liquidityPoolDepositOp

		return result, nil
	case xdr.OperationTypeLiquidityPoolWithdraw:
		xdrLiquidityPoolWithdrawOp := bd.LiquidityPoolWithdrawOp

		xdrHashLpId := xdr.Hash(xdrLiquidityPoolWithdrawOp.LiquidityPoolId)
		LpId := PoolId(xdrHashLpId[:])

		liquidityPoolWithdrawOp := &LiquidityPoolWithdrawOp{
			LiquidityPoolId: LpId,
			Amount:          int64(xdrLiquidityPoolWithdrawOp.Amount),
			MinAmountA:      int64(xdrLiquidityPoolWithdrawOp.MinAmountA),
			MinAmountB:      int64(xdrLiquidityPoolWithdrawOp.MinAmountB),
		}
		result.LiquidityPoolWithdrawOp = liquidityPoolWithdrawOp

		return result, nil
	case xdr.OperationTypeInvokeHostFunction:
		xdrInvokeHostFunctionOp := bd.InvokeHostFunctionOp

		invokeHostFunctionOp := &InvokeHostFunctionOp{}
	case xdr.OperationTypeExtendFootprintTtl:
	case xdr.OperationTypeRestoreFootprint:
	}
	return OperationBody{}, nil
}
