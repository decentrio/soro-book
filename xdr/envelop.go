package xdr

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertEnvelopeXdrToJson(envelope xdr.TransactionEnvelope) (TransactionEnvelope, error) {
	// var txEnvelope TransactionEnvelope
	switch envelope.Type {
	case xdr.EnvelopeTypeEnvelopeTypeTxV0:
	case xdr.EnvelopeTypeEnvelopeTypeTx:
	case xdr.EnvelopeTypeEnvelopeTypeTxFeeBump:
	}

	return TransactionEnvelope{}, errors.Errorf("error invalid type envelope: %v", envelope.Type)
}

func ConvertEnvelopeV0XdrToJson(v0 *xdr.TransactionV0Envelope) (TransactionV0Envelope, error) {
	return TransactionV0Envelope{}, nil
}
