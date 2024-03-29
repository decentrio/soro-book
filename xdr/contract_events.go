package xdr

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

func ConvertContractEvent(e xdr.ContractEvent) ContractEvent {

}

// func ConvertContractEventJSON(event models.Event, topics []models.Topics) (*models.EventJSON, error) {
// 	evt := &models.EventJSON{}

// 	evt.Id = event.Id
// 	evt.ContractId = event.ContractId
// 	evt.TxHash = event.TxHash
// 	evt.EventType = event.EventType

// 	var xdrTopics []xdr.ScVal
// 	for _, topic := range topics {
// 		var xdrTopic xdr.ScVal
// 		err := xdrTopic.UnmarshalBinary([]byte(topic.TopicXdr))
// 		if err != nil {
// 			return evt, fmt.Errorf("Error Unmarshal topic binary")
// 		}
// 		xdrTopics = append(xdrTopics, xdrTopic)
// 	}

// 	var value xdr.ScVal
// 	err := value.UnmarshalBinary([]byte(event.ValueXdr))
// 	if err != nil {
// 		return evt, fmt.Errorf("Error Unmarshal value binary")
// 	}

// 	switch evt.EventType {
// 	case aggregation.EventTypeTransfer:
// 		transferEvent := TransferEvent{}
// 		transferEvent.parse(xdrTopics, value)

// 		bz, err := transferEvent.ToJSON()
// 		if err != nil {
// 			return evt, err
// 		}

// 		evt.Data = bz
// 	case aggregation.EventTypeMint:
// 		mintEvent := MintEvent{}
// 		mintEvent.parse(xdrTopics, value)

// 		bz, err := mintEvent.ToJSON()
// 		if err != nil {
// 			return evt, err
// 		}

// 		evt.Data = bz
// 	case aggregation.EventTypeClawback:
// 		cbEvent := ClawbackEvent{}
// 		cbEvent.parse(xdrTopics, value)

// 		bz, err := cbEvent.ToJSON()
// 		if err != nil {
// 			return evt, err
// 		}

// 		evt.Data = bz
// 	case aggregation.EventTypeBurn:
// 		burnEvent := BurnEvent{}
// 		burnEvent.parse(xdrTopics, value)

// 		bz, err := burnEvent.ToJSON()
// 		if err != nil {
// 			return evt, err
// 		}

// 		evt.Data = bz
// 	default:
// 		return evt, errors.Wrapf(aggregation.ErrEventUnsupported, "event not supported %s", evt.EventType)
// 	}

// 	return evt, nil
// }

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
