package converter

import (
	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

func ConvertPreconditions(c xdr.Preconditions) (Preconditions, error) {
	var result Preconditions
	switch c.Type {
	case xdr.PreconditionTypePrecondNone:
		return result, nil
	case xdr.PreconditionTypePrecondTime:
		timeBounds, err := ConvertTimeBounds(c.TimeBounds)
		if err != nil {
			return result, err
		}
		result.TimeBounds = timeBounds

		return result, nil
	case xdr.PreconditionTypePrecondV2:
		v2, err := ConvertPreconditionsV2(*c.V2)
		if err != nil {
			return result, err
		}
		result.V2 = &v2

		return result, nil
	}

	return result, errors.Errorf("error invalid Preconditions type %v", c.Type)
}

func ConvertPreconditionsV2(c xdr.PreconditionsV2) (PreconditionsV2, error) {
	var result PreconditionsV2

	var timeBounds *TimeBounds
	var ledgerBounds LedgerBounds

	timeBounds, err := ConvertTimeBounds(c.TimeBounds)
	if err != nil {
		return result, err
	}

	if c.LedgerBounds != nil {
		ledgerBounds = ConvertLedgerBounds(*c.LedgerBounds)
	}

	var minSeqNum int64
	if c.MinSeqNum != nil {
		minSeqNum = int64(*c.MinSeqNum)
	}

	var extraSigners []SignerKey
	for _, xdrSigner := range c.ExtraSigners {
		signer, err := ConvertSignerKey(xdrSigner)
		if err != nil {
			return result, err
		}
		extraSigners = append(extraSigners, signer)
	}

	result.TimeBounds = timeBounds
	result.LedgerBounds = &ledgerBounds
	result.MinSeqNum = &minSeqNum
	result.MinSeqAge = uint64(c.MinSeqAge)
	result.MinSeqLedgerGap = uint32(c.MinSeqLedgerGap)
	result.ExtraSigners = extraSigners

	return result, nil
}
