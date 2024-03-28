package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

// TODO: testing
func ConvertTransactionEnvelope(e xdr.TransactionEnvelope) (TransactionEnvelope, error) {
	var result TransactionEnvelope
	switch e.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
		v0, err := ConvertTransactionV0Envelope(e.V0)
		if err != nil {
			return result, err
		}
		result.V0 = &v0

		return result, nil
	case xdr.EnvelopeTypeEnvelopeTypeTx:
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
	}

	return TransactionEnvelope{}, errors.Errorf("error invalid type envelope: %v", e.Type)
}

// TODO: testing
func ConvertTransactionV0Envelope(v0 *xdr.TransactionV0Envelope) (TransactionV0Envelope, error) {
	var result TransactionV0Envelope
	tx, err := ConvertTransactionV0(v0.Tx)
	if err != nil {
		return result, err
	}

	var sigs []DecoratedSignature
	for _, xdrSig := range v0.Signatures {
		sig := ConvertDecoratedSignature(xdrSig)
		sigs = append(sigs, sig)
	}

	result.Tx = tx
	result.Signatures = sigs

	return result, nil
}

// TODO: testing
func ConvertTransactionV1Envelope(v1 *xdr.TransactionV1Envelope) (TransactionV1Envelope, error) {
	var result TransactionV1Envelope
	tx, err := ConvertTransaction(v1.Tx)
	if err != nil {
		return result, err
	}

	var sigs []DecoratedSignature
	for _, xdrSig := range v1.Signatures {
		sig := ConvertDecoratedSignature(xdrSig)
		sigs = append(sigs, sig)
	}

	result.Tx = tx
	result.Signatures = sigs

	return result, nil
}
