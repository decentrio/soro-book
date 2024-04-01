package converter

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertOperationMeta(m xdr.OperationMeta) (OperationMeta, error) {
	var result OperationMeta
	var changes LedgerEntryChanges
	for _, xdrChange := range m.Changes {
		change, err := ConvertLedgerEntryChange(xdrChange)
		if err != nil {
			return result, err
		}
		changes = append(changes, change)
	}
	result.Changes = changes

	return result, nil
}

func ConvertOperationResult(op xdr.OperationResult) (OperationResult, error) {
	var result OperationResult
	result.Code = int32(op.Code)

	if op.Code == xdr.OperationResultCodeOpInner {
		tr, err := ConvertOperationResultTr(*op.Tr)
		if err != nil {
			return result, err
		}

		result.Tr = &tr
	}

	return result, nil
}

func ConvertOperationResultTr(r xdr.OperationResultTr) (OperationResultTr, error) {
	var result OperationResultTr

	switch r.Type {
	case xdr.OperationTypeCreateAccount:
		xdrCreateAccountResult := r.CreateAccountResult

		createAccountResult := CreateAccountResult{
			Code: int32(xdrCreateAccountResult.Code),
		}
		result.CreateAccountResult = &createAccountResult

		return result, nil
	case xdr.OperationTypePayment:
		xdrPaymentResult := r.PaymentResult

		paymentResult := PaymentResult{
			Code: int32(xdrPaymentResult.Code),
		}
		result.PaymentResult = &paymentResult

		return result, nil
	case xdr.OperationTypePathPaymentStrictReceive:
		xdrPathPaymentStrictReceiveResult := r.PathPaymentStrictReceiveResult

		pathPaymentStrictReceiveResult := PathPaymentStrictReceiveResult{
			Code: int32(xdrPathPaymentStrictReceiveResult.Code),
		}

		if xdrPathPaymentStrictReceiveResult.Code == xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveSuccess {
			success, err := ConvertPathPaymentStrictReceiveResultSuccess(*xdrPathPaymentStrictReceiveResult.Success)
			if err != nil {
				return result, err
			}
			pathPaymentStrictReceiveResult.Success = &success
		} else if xdrPathPaymentStrictReceiveResult.Code == xdr.PathPaymentStrictReceiveResultCodePathPaymentStrictReceiveNoIssuer {

			noIssuer, err := ConvertAsset(*xdrPathPaymentStrictReceiveResult.NoIssuer)
			if err != nil {
				return result, err
			}
			pathPaymentStrictReceiveResult.NoIssuer = &noIssuer
		}

		result.PathPaymentStrictReceiveResult = &pathPaymentStrictReceiveResult

		return result, nil
	case xdr.OperationTypeManageSellOffer:
		xdrManageSellOfferResult := r.ManageSellOfferResult

		manageSellOfferResult := ManageSellOfferResult{
			Code: int32(xdrManageSellOfferResult.Code),
		}

		if xdrManageSellOfferResult.Code == xdr.ManageSellOfferResultCodeManageSellOfferSuccess {
			success, err := ConvertManageOfferSuccessResult(*xdrManageSellOfferResult.Success)
			if err != nil {
				return result, err
			}

			manageSellOfferResult.Success = &success
		}
		result.ManageSellOfferResult = &manageSellOfferResult

		return result, nil
	case xdr.OperationTypeCreatePassiveSellOffer:
		xdrCreatePassiveSellOfferResult := r.CreatePassiveSellOfferResult

		createPassiveSellOfferResult := ManageSellOfferResult{
			Code: int32(xdrCreatePassiveSellOfferResult.Code),
		}

		if xdrCreatePassiveSellOfferResult.Code == xdr.ManageSellOfferResultCodeManageSellOfferSuccess {
			success, err := ConvertManageOfferSuccessResult(*xdrCreatePassiveSellOfferResult.Success)
			if err != nil {
				return result, err
			}

			createPassiveSellOfferResult.Success = &success
		}
		result.ManageSellOfferResult = &createPassiveSellOfferResult
	case xdr.OperationTypeSetOptions:
		xdrSetOptionsResult := r.SetOptionsResult

		setOptionsResult := SetOptionsResult{
			Code: int32(xdrSetOptionsResult.Code),
		}
		result.SetOptionsResult = &setOptionsResult

		return result, nil
	case xdr.OperationTypeChangeTrust:
		xdrChangeTrustResult := r.ChangeTrustResult

		changeTrustResult := ChangeTrustResult{
			Code: int32(xdrChangeTrustResult.Code),
		}
		result.ChangeTrustResult = &changeTrustResult

		return result, nil
	case xdr.OperationTypeAllowTrust:
		xdrAllowTrustResult := r.AllowTrustResult

		allowTrustResult := AllowTrustResult{
			Code: int32(xdrAllowTrustResult.Code),
		}
		result.AllowTrustResult = &allowTrustResult

		return result, nil
	case xdr.OperationTypeAccountMerge:
		xdrAccountMergeResult := r.AccountMergeResult

		accountMergeResult := AccountMergeResult{
			Code: int32(xdrAccountMergeResult.Code),
		}

		if xdrAccountMergeResult.Code == xdr.AccountMergeResultCodeAccountMergeSuccess {
			sourceAccountBalance := int64(*xdrAccountMergeResult.SourceAccountBalance)
			accountMergeResult.SourceAccountBalance = &sourceAccountBalance
		}
		result.AccountMergeResult = &accountMergeResult

		return result, nil
	case xdr.OperationTypeInflation:
		xdrInflationResult := r.InflationResult

		inflationResult := InflationResult{
			Code: int32(xdrInflationResult.Code),
		}

		if xdrInflationResult.Code == xdr.InflationResultCodeInflationSuccess {
			var payouts []InflationPayout
			for _, xdrPayout := range *xdrInflationResult.Payouts {
				payout := ConvertInflationPayout(xdrPayout)
				payouts = append(payouts, payout)
			}
			inflationResult.Payouts = &payouts
		}
		result.InflationResult = &inflationResult

		return result, nil
	case xdr.OperationTypeManageData:
		xdrManageDataResult := r.ManageDataResult

		manageDataResult := ManageDataResult{
			Code: int32(xdrManageDataResult.Code),
		}
		result.ManageDataResult = &manageDataResult

		return result, nil
	case xdr.OperationTypeBumpSequence:
		xdrBumpSeqResult := r.BumpSeqResult

		bumpSequenceResult := BumpSequenceResult{
			Code: int32(xdrBumpSeqResult.Code),
		}
		result.BumpSeqResult = &bumpSequenceResult

		return result, nil
	case xdr.OperationTypeManageBuyOffer:
		xdrManageBuyOfferResult := r.ManageBuyOfferResult

		manageBuyOfferResult := ManageBuyOfferResult{
			Code: int32(xdrManageBuyOfferResult.Code),
		}

		if xdrManageBuyOfferResult.Code == xdr.ManageBuyOfferResultCodeManageBuyOfferSuccess {
			success, err := ConvertManageOfferSuccessResult(*xdrManageBuyOfferResult.Success)
			if err != nil {
				return result, err
			}

			manageBuyOfferResult.Success = &success
		}
		result.ManageBuyOfferResult = &manageBuyOfferResult

		return result, nil
	case xdr.OperationTypePathPaymentStrictSend:
		xdrPathPaymentStrictSendResult := r.PathPaymentStrictSendResult

		pathPaymentStrictSendResult := PathPaymentStrictSendResult{
			Code: int32(xdrPathPaymentStrictSendResult.Code),
		}

		if xdrPathPaymentStrictSendResult.Code == xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendSuccess {
			success, err := ConvertPathPaymentStrictSendResultSuccess(*xdrPathPaymentStrictSendResult.Success)
			if err != nil {
				return result, err
			}
			pathPaymentStrictSendResult.Success = &success
		} else if xdrPathPaymentStrictSendResult.Code == xdr.PathPaymentStrictSendResultCodePathPaymentStrictSendNoIssuer {

			noIssuer, err := ConvertAsset(*xdrPathPaymentStrictSendResult.NoIssuer)
			if err != nil {
				return result, err
			}
			pathPaymentStrictSendResult.NoIssuer = &noIssuer
		}
		result.PathPaymentStrictSendResult = &pathPaymentStrictSendResult

		return result, nil
	case xdr.OperationTypeCreateClaimableBalance:
		xdrCreateClaimableBalanceResult := r.CreateClaimableBalanceResult

		createClaimableBalanceResult := CreateClaimableBalanceResult{
			Code: int32(xdrCreateClaimableBalanceResult.Code),
		}

		if xdrCreateClaimableBalanceResult.Code == xdr.CreateClaimableBalanceResultCodeCreateClaimableBalanceSuccess {
			balanceId, err := ConvertClaimableBalanceId(*xdrCreateClaimableBalanceResult.BalanceId)
			if err != nil {
				return result, err
			}
			createClaimableBalanceResult.BalanceId = &balanceId
		}
		result.CreateClaimableBalanceResult = &createClaimableBalanceResult

		return result, nil
	case xdr.OperationTypeClaimClaimableBalance:
		xdrClaimClaimableBalanceResult := r.ClaimClaimableBalanceResult

		claimClaimableBalanceResult := ClaimClaimableBalanceResult{
			Code: int32(xdrClaimClaimableBalanceResult.Code),
		}
		result.ClaimClaimableBalanceResult = &claimClaimableBalanceResult

		return result, nil
	case xdr.OperationTypeBeginSponsoringFutureReserves:
		xdrBeginSponsoringFutureReservesResult := r.BeginSponsoringFutureReservesResult

		beginSponsoringFutureReservesResult := BeginSponsoringFutureReservesResult{
			Code: int32(xdrBeginSponsoringFutureReservesResult.Code),
		}
		result.BeginSponsoringFutureReservesResult = &beginSponsoringFutureReservesResult

		return result, nil
	case xdr.OperationTypeEndSponsoringFutureReserves:
		xdrEndSponsoringFutureReservesResult := r.EndSponsoringFutureReservesResult

		endSponsoringFutureReservesResult := EndSponsoringFutureReservesResult{
			Code: int32(xdrEndSponsoringFutureReservesResult.Code),
		}
		result.EndSponsoringFutureReservesResult = &endSponsoringFutureReservesResult

		return result, nil
	case xdr.OperationTypeRevokeSponsorship:
		xdrRevokeSponsorshipResult := r.RevokeSponsorshipResult

		revokeSponsorshipResult := RevokeSponsorshipResult{
			Code: int32(xdrRevokeSponsorshipResult.Code),
		}
		result.RevokeSponsorshipResult = &revokeSponsorshipResult

		return result, nil
	case xdr.OperationTypeClawback:
		xdrClawbackResult := r.ClawbackResult

		clawbackResult := ClawbackResult{
			Code: int32(xdrClawbackResult.Code),
		}
		result.ClawbackResult = &clawbackResult

		return result, nil
	case xdr.OperationTypeClawbackClaimableBalance:
		xdrClawbackClaimableBalanceResult := r.ClawbackClaimableBalanceResult

		clawbackClaimableBalanceResult := ClawbackClaimableBalanceResult{
			Code: int32(xdrClawbackClaimableBalanceResult.Code),
		}
		result.ClawbackClaimableBalanceResult = &clawbackClaimableBalanceResult

		return result, nil
	case xdr.OperationTypeSetTrustLineFlags:
		xdrSetTrustLineFlagsResult := r.SetTrustLineFlagsResult

		setTrustLineFlagsResult := SetTrustLineFlagsResult{
			Code: int32(xdrSetTrustLineFlagsResult.Code),
		}
		result.SetTrustLineFlagsResult = &setTrustLineFlagsResult

		return result, nil
	case xdr.OperationTypeLiquidityPoolDeposit:
		xdrLiquidityPoolDepositResult := r.LiquidityPoolDepositResult

		liquidityPoolDepositResult := LiquidityPoolDepositResult{
			Code: int32(xdrLiquidityPoolDepositResult.Code),
		}
		result.LiquidityPoolDepositResult = &liquidityPoolDepositResult

		return result, nil
	case xdr.OperationTypeLiquidityPoolWithdraw:
		xdrLiquidityPoolWithdrawResult := r.LiquidityPoolWithdrawResult

		liquidityPoolWithdrawResult := LiquidityPoolWithdrawResult{
			Code: int32(xdrLiquidityPoolWithdrawResult.Code),
		}
		result.LiquidityPoolWithdrawResult = &liquidityPoolWithdrawResult

		return result, nil
	case xdr.OperationTypeInvokeHostFunction:
		xdrInvokeHostFunctionResult := r.InvokeHostFunctionResult

		invokeHostFunctionResult := InvokeHostFunctionResult{
			Code: int32(xdrInvokeHostFunctionResult.Code),
		}

		if xdrInvokeHostFunctionResult.Code == xdr.InvokeHostFunctionResultCodeInvokeHostFunctionSuccess {
			success := (*xdrInvokeHostFunctionResult.Success).HexString()
			invokeHostFunctionResult.Success = &success
		}
		result.InvokeHostFunctionResult = &invokeHostFunctionResult

		return result, nil
	case xdr.OperationTypeExtendFootprintTtl:
		xdrExtendFootprintTtlResult := r.ExtendFootprintTtlResult

		extendFootprintTtlResult := ExtendFootprintTtlResult{
			Code: int32(xdrExtendFootprintTtlResult.Code),
		}
		result.ExtendFootprintTtlResult = &extendFootprintTtlResult

		return result, nil
	case xdr.OperationTypeRestoreFootprint:
		xdrRestoreFootprintResult := r.RestoreFootprintResult

		restoreFootprintResult := RestoreFootprintResult{
			Code: int32(xdrRestoreFootprintResult.Code),
		}
		result.RestoreFootprintResult = &restoreFootprintResult

		return result, nil
	}

	return result, errors.Errorf("error invalid operationBody key type %v", r.Type)
}

// TODO: testing
func ConvertOperation(op xdr.Operation) (Operation, error) {
	var result Operation
	var sourceAccount MuxedAccount
	var err error
	if op.SourceAccount != nil {
		sourceAccount, err = ConvertMuxedAccount(*op.SourceAccount)
		if err != nil {
			return result, err
		}
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
		xdrManageSellOffer := bd.ManageSellOfferOp

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
			BuyAmount: int64(xdrManageSellOffer.Amount),
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

		var inflationDest PublicKey
		var clearFlags, setFlags, masterWeight, lowThreshold, medThreshold, highThreshold uint32
		var homeDomain string
		var signer Signer

		if xdrSetOptions.InflationDest != nil {
			inflationDest = PublicKey{
				Ed25519: ConvertEd25519(xdrSetOptions.InflationDest.Ed25519),
			}
		}

		if xdrSetOptions.ClearFlags != nil {
			clearFlags = uint32(*xdrSetOptions.ClearFlags)
		}

		if xdrSetOptions.SetFlags != nil {
			setFlags = uint32(*xdrSetOptions.SetFlags)
		}

		if xdrSetOptions.MasterWeight != nil {
			masterWeight = uint32(*xdrSetOptions.MasterWeight)
		}

		if xdrSetOptions.LowThreshold != nil {
			lowThreshold = uint32(*xdrSetOptions.LowThreshold)
		}

		if xdrSetOptions.MedThreshold != nil {
			medThreshold = uint32(*xdrSetOptions.MedThreshold)
		}

		if xdrSetOptions.HighThreshold != nil {
			highThreshold = uint32(*xdrSetOptions.HighThreshold)
		}

		if xdrSetOptions.HomeDomain != nil {
			homeDomain = string(*xdrSetOptions.HomeDomain)
		}

		var err error
		if xdrSetOptions.Signer != nil {
			signer, err = ConvertSigner(*xdrSetOptions.Signer)
			if err != nil {
				return result, err
			}
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
			return result, err
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
			return result, err
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

		var ledgerKey LedgerKey
		var signer RevokeSponsorshipOpSigner
		var err error

		if xdrRevokeSponsorshipOp.LedgerKey != nil {
			ledgerKey, err = ConvertLedgerKey(*xdrRevokeSponsorshipOp.LedgerKey)
			if err != nil {
				return result, err
			}
		}

		if xdrRevokeSponsorshipOp.Signer != nil {
			signer, err = ConvertRevokeSponsorshipOpSigner(*xdrRevokeSponsorshipOp.Signer)
			if err != nil {
				return result, err
			}
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
			return result, err
		}

		from, err := ConvertMuxedAccount(xdrClawbackOp.From)
		if err != nil {
			return result, err
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
			return result, err
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
			return result, err
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

		hostFunc, err := ConvertHostFunction(xdrInvokeHostFunctionOp.HostFunction)
		if err != nil {
			return result, err
		}

		var auths []SorobanAuthorizationEntry
		for _, xdrEntry := range xdrInvokeHostFunctionOp.Auth {
			auth, err := ConvertSorobanAuthorizationEntry(xdrEntry)
			if err != nil {
				return result, err
			}

			auths = append(auths, auth)
		}

		invokeHostFunctionOp := &InvokeHostFunctionOp{
			HostFunction: hostFunc,
			Auth:         auths,
		}
		result.InvokeHostFunctionOp = invokeHostFunctionOp

		return result, nil
	case xdr.OperationTypeExtendFootprintTtl:
		xdrExtendFootprintTtlOp := bd.ExtendFootprintTtlOp

		extendFootprintTtlOp := &ExtendFootprintTtlOp{
			Ext:      ConvertExtensionPoint(xdrExtendFootprintTtlOp.Ext),
			ExtendTo: uint32(xdrExtendFootprintTtlOp.ExtendTo),
		}
		result.ExtendFootprintTtlOp = extendFootprintTtlOp

		return result, nil
	case xdr.OperationTypeRestoreFootprint:
		xdrRestoreFootprintOp := bd.RestoreFootprintOp

		restoreFootprintOp := &RestoreFootprintOp{
			Ext: ConvertExtensionPoint(xdrRestoreFootprintOp.Ext),
		}
		result.RestoreFootprintOp = restoreFootprintOp

		return result, nil
	}
	return result, errors.Errorf("error invalid operationBody key type %v", bd.Type)
}

func ConvertStateArchivalSettings(s xdr.StateArchivalSettings) StateArchivalSettings {
	return StateArchivalSettings{
		MaxEntryTtl:                    uint32(s.MaxEntryTtl),
		MinTemporaryTtl:                uint32(s.MinTemporaryTtl),
		MinPersistentTtl:               uint32(s.MinPersistentTtl),
		PersistentRentRateDenominator:  int64(s.PersistentRentRateDenominator),
		TempRentRateDenominator:        int64(s.TempRentRateDenominator),
		MaxEntriesToArchive:            uint32(s.MaxEntriesToArchive),
		BucketListSizeWindowSampleSize: uint32(s.BucketListSizeWindowSampleSize),
		BucketListWindowSamplePeriod:   uint32(s.BucketListWindowSamplePeriod),
		EvictionScanSize:               uint32(s.EvictionScanSize),
		StartingEvictionScanLevel:      uint32(s.StartingEvictionScanLevel),
	}
}
