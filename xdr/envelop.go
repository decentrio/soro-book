package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

// TODO: testing
func ConvertTransactionEnvelope(envelope xdr.TransactionEnvelope) (TransactionEnvelope, error) {
	// var txEnvelope TransactionEnvelope
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:

	case xdr.EnvelopeTypeEnvelopeTypeTx:
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
	}

	return TransactionEnvelope{}, errors.Errorf("error invalid type envelope: %v", envelope.Type)
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

	return TransactionV0Envelope{
		Tx:         tx,
		Signatures: sigs,
	}, nil
}
