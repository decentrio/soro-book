package aggregation

import (
	"github.com/pkg/errors"

	"github.com/stellar/go/xdr"
)

const (
	// Implemented
	EventTypeTransfer = "transfer"
	EventTypeMint     = "mint"
	EventTypeClawback = "clawback"
	EventTypeBurn     = "burn"
	// TODO: Not implemented
	EventTypeIncrAllow
	EventTypeDecrAllow
	EventTypeSetAuthorized
	EventTypeSetAdmin
)

var (
	STELLAR_ASSET_CONTRACT_TOPICS = map[xdr.ScSymbol]string{
		xdr.ScSymbol("transfer"): EventTypeTransfer,
		xdr.ScSymbol("mint"):     EventTypeMint,
		xdr.ScSymbol("clawback"): EventTypeClawback,
		xdr.ScSymbol("burn"):     EventTypeBurn,
	}

	ErrNotStellarAssetContract = errors.New("event was not from a Stellar Asset Contract")
	ErrEventUnsupported        = errors.New("this type of Stellar Asset Contract event is unsupported")
	ErrEventIntegrity          = errors.New("contract ID doesn't match asset + passphrase")
)

func getEventType(eventBody xdr.ContractEventBody) (string, bool) {
	topics := eventBody.V0.Topics

	if len(topics) <= 2 {
		return "", false
	}

	// Filter out events for function calls we don't care about
	fn, ok := topics[0].GetSym()
	if !ok {
		return "", false
	}

	eventType, found := STELLAR_ASSET_CONTRACT_TOPICS[fn]
	if !found {
		return "", false
	}

	return eventType, true
}
