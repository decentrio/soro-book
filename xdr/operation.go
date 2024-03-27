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

	case xdr.OperationTypeManageSellOffer:
	case xdr.OperationTypeCreatePassiveSellOffer:
	case xdr.OperationTypeSetOptions:
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
