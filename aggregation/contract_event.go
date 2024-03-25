package aggregation

import (
	"encoding/json"
	"fmt"

	"github.com/decentrio/soro-book/database/models"
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

	ErrNotTransferEvent = errors.New("this is not transfer event")
	ErrNotMintEvent     = errors.New("this is not mint event")
	ErrNotClawbackEvent = errors.New("this is not clawback event")
	ErrNotBurnEvent     = errors.New("this is not burn event")
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

func ContractEventJSON(event models.ContractEvent, topics []models.Topics) (*models.ContractEventJSON, error) {
	evt := &models.ContractEventJSON{}

	evt.Id = event.Id
	evt.ContractId = event.ContractId
	evt.LedgerSeq = event.LedgerSeq
	evt.TxHash = event.TxHash
	evt.EventType = event.EventType

	var xdrTopics []xdr.ScVal
	for _, topic := range topics {
		var xdrTopic xdr.ScVal
		err := xdrTopic.UnmarshalBinary([]byte(topic.TopicXdr))
		if err != nil {
			return evt, fmt.Errorf("Error Unmarshal topic binary")
		}
		xdrTopics = append(xdrTopics, xdrTopic)
	}

	var value xdr.ScVal
	err := value.UnmarshalBinary([]byte(event.Value))
	if err != nil {
		return evt, fmt.Errorf("Error Unmarshal value binary")
	}

	switch evt.EventType {
	case EventTypeTransfer:
		transferEvent := TransferEvent{}
		transferEvent.parse(xdrTopics, value)

		bz, err := transferEvent.ToJSON()
		if err != nil {
			return evt, err
		}

		evt.Data = bz
	case EventTypeMint:
		mintEvent := MintEvent{}
		mintEvent.parse(xdrTopics, value)

		bz, err := mintEvent.ToJSON()
		if err != nil {
			return evt, err
		}

		evt.Data = bz
	case EventTypeClawback:
		cbEvent := ClawbackEvent{}
		cbEvent.parse(xdrTopics, value)

		bz, err := cbEvent.ToJSON()
		if err != nil {
			return evt, err
		}

		evt.Data = bz
	case EventTypeBurn:
		burnEvent := BurnEvent{}
		burnEvent.parse(xdrTopics, value)

		bz, err := burnEvent.ToJSON()
		if err != nil {
			return evt, err
		}

		evt.Data = bz
	default:
		return evt, errors.Wrapf(ErrEventUnsupported, "event not supported %s", evt.EventType)
	}

	return evt, nil
}

type Int128Parts struct {
	Hi int64  `json:"hi,omitempty"`
	Lo uint64 `json:"lo,omitempty"`
}

type TransferEvent struct {
	From   string      `json:"from,omitempty"`
	To     string      `json:"to,omitempty"`
	Amount Int128Parts `json:"amount,omitempty"`
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

type MintEvent struct {
	Admin  string      `json:"admin,omitempty"`
	To     string      `json:"to,omitempty"`
	Amount Int128Parts `json:"amount,omitempty"`
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

type ClawbackEvent struct {
	Admin  string      `json:"admin,omitempty"`
	From   string      `json:"from,omitempty"`
	Amount Int128Parts `json:"amount,omitempty"`
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

type BurnEvent struct {
	From   string      `json:"from,omitempty"`
	Amount Int128Parts `json:"amount,omitempty"`
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
