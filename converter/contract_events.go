package converter

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

var (
	ErrNotBalanceChangeEvent = errors.New("event doesn't represent a balance change")
	ErrNotTransferEvent      = errors.New("this is not transfer event")
	ErrNotMintEvent          = errors.New("this is not mint event")
	ErrNotClawbackEvent      = errors.New("this is not clawback event")
	ErrNotBurnEvent          = errors.New("this is not burn event")
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
		topics[0].GetSym()

		return string(fn), false
	}

	return eventType, true
}

func ConvertContractEvent(e xdr.ContractEvent) (ContractEvent, error) {
	var result ContractEvent

	result.Ext = ConvertExtensionPoint(e.Ext)
	contractId := e.ContractId.HexString()
	result.ContractId = &contractId
	result.ContractEventType = int32(e.Type)

	eventType, found := getEventType(e.Body)
	result.EventType = eventType
	if !found {
		return result, nil
	}

	topics := e.Body.V0.Topics
	value := e.Body.V0.Data

	switch eventType {
	case EventTypeTransfer:
		transferEvent := TransferEvent{}
		transferEvent.parse(topics, value)
		result.Transfer = &transferEvent
	case EventTypeMint:
		mintEvent := MintEvent{}
		mintEvent.parse(topics, value)
		result.Mint = &mintEvent
	case EventTypeClawback:
		cbEvent := ClawbackEvent{}
		cbEvent.parse(topics, value)
		result.Clawback = &cbEvent
	case EventTypeBurn:
		burnEvent := BurnEvent{}
		burnEvent.parse(topics, value)
		result.Burn = &burnEvent
	default:
		return result, errors.Wrapf(ErrEventUnsupported, "event not supported %s", eventType)
	}
	return result, nil
}

func ConvertDiagnosticEvent(e xdr.DiagnosticEvent) (DiagnosticEvent, error) {
	var result DiagnosticEvent

	event, err := ConvertContractEvent(e.Event)
	if err != nil {
		return result, err
	}

	result.InSuccessfulContractCall = e.InSuccessfulContractCall
	result.Event = event

	return result, nil
}

func (event *TransferEvent) parse(topics xdr.ScVec, value xdr.ScVal) error {
	//
	// The transfer event format is:
	//
	// 	"transfer"  Symbol
	//  <from> 		Address
	//  <to> 		Address
	// 	<asset>		Bytes
	//
	// 	<amount> 	i128
	//
	var err error
	event.From, event.To, event.Amount, err = parseBalanceChangeEvent(topics, value)
	if err != nil {
		return ErrNotTransferEvent
	}
	return nil
}

func (e TransferEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (event *MintEvent) parse(topics xdr.ScVec, value xdr.ScVal) error {
	//
	// The mint event format is:
	//
	// 	"mint"  	Symbol
	//  <admin>		Address
	//  <to> 		Address
	// 	<asset>		Bytes
	//
	// 	<amount> 	i128
	//
	var err error
	event.Admin, event.To, event.Amount, err = parseBalanceChangeEvent(topics, value)
	if err != nil {
		return ErrNotMintEvent
	}
	return nil
}

func (e MintEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (event *ClawbackEvent) parse(topics xdr.ScVec, value xdr.ScVal) error {
	//
	// The clawback event format is:
	//
	// 	"clawback" 	Symbol
	//  <admin>		Address
	//  <from> 		Address
	// 	<asset>		Bytes
	//
	// 	<amount> 	i128
	//
	var err error
	event.Admin, event.From, event.Amount, err = parseBalanceChangeEvent(topics, value)
	if err != nil {
		return ErrNotClawbackEvent
	}
	return nil
}

func (e ClawbackEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (event *BurnEvent) parse(topics xdr.ScVec, value xdr.ScVal) error {
	//
	// The transfer event format is:
	//
	// 	"burn"  	Symbol
	//  <from> 		Address
	//  <to> 		Address
	// 	<asset>		Bytes
	//
	// 	<amount> 	i128
	//
	if len(topics) != 3 {
		return ErrNotBurnEvent
	}

	from, ok := topics[1].GetAddress()
	if !ok {
		return ErrNotBurnEvent
	}

	var err error
	event.From, err = from.String()
	if err != nil {
		return errors.Wrap(err, ErrNotBurnEvent.Error())
	}

	amount, ok := value.GetI128()
	if !ok {
		return ErrNotBurnEvent
	}
	event.Amount = XdrInt128PartsConvert(amount)

	return nil
}

func (e BurnEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

// parseBalanceChangeEvent is a generalization of a subset of the Stellar Asset
// Contract events. Transfer, mint, clawback, and burn events all have two
// addresses and an amount involved. The addresses represent different things in
// different event types (e.g. "from" or "admin"), but the parsing is identical.
// This helper extracts all three parts or returns a generic error if it can't.
func parseBalanceChangeEvent(topics xdr.ScVec, value xdr.ScVal) (
	first string,
	second string,
	amount Int128Parts,
	err error,
) {
	err = ErrNotBalanceChangeEvent
	if len(topics) != 4 {
		return
	}

	firstSc, ok := topics[1].GetAddress()
	if !ok {
		return
	}
	first, err = firstSc.String()
	if err != nil {
		err = errors.Wrap(err, ErrNotBalanceChangeEvent.Error())
		return
	}

	secondSc, ok := topics[2].GetAddress()
	if !ok {
		return
	}
	second, err = secondSc.String()
	if err != nil {
		err = errors.Wrap(err, ErrNotBalanceChangeEvent.Error())
		return
	}

	xdrAmount, ok := value.GetI128()
	if !ok {
		return
	}

	amount = XdrInt128PartsConvert(xdrAmount)

	return first, second, amount, nil
}

func XdrInt128PartsConvert(in xdr.Int128Parts) Int128Parts {
	out := Int128Parts{
		Hi: int64(in.Hi),
		Lo: uint64(in.Lo),
	}

	return out
}
