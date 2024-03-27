package xdr

import (
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
	case xdr.OperationTypeAllowTrust:
	case xdr.OperationTypeAccountMerge:
	case xdr.OperationTypeInflation:
	case xdr.OperationTypeManageData:
	case xdr.OperationTypeBumpSequence:
	case xdr.OperationTypeManageBuyOffer:
	case xdr.OperationTypePathPaymentStrictSend:
	case xdr.OperationTypeCreateClaimableBalance:
	case xdr.OperationTypeClaimClaimableBalance:
	case xdr.OperationTypeBeginSponsoringFutureReserves:
	case xdr.OperationTypeEndSponsoringFutureReserves:
	case xdr.OperationTypeRevokeSponsorship:
	case xdr.OperationTypeClawback:
	case xdr.OperationTypeClawbackClaimableBalance:
	case xdr.OperationTypeSetTrustLineFlags:
	case xdr.OperationTypeLiquidityPoolDeposit:
	case xdr.OperationTypeLiquidityPoolWithdraw:
	case xdr.OperationTypeInvokeHostFunction:
	case xdr.OperationTypeExtendFootprintTtl:
	case xdr.OperationTypeRestoreFootprint:
	}
	return OperationBody{}, nil
}
