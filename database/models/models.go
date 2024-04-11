package models

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/stellar/go/xdr"
)

type Ledger struct {
	Hash         string `json:"hash,omitempty"`
	PrevHash     string `json:"prev_hash,omitempty"`
	Sequence     uint32 `json:"sequence,omitempty"`
	Transactions uint32 `json:"transaction,omitempty"`
	Operations   uint32 `json:"operations,omitempty"`
}

type Transaction struct {
	Hash             string `json:"hash,omitempty"`
	Status           string `json:"status,omitempty"`
	Ledger           uint32 `json:"ledger,omitempty"`
	ApplicationOrder uint32 `json:"application_order,omitempty"`
	EnvelopeXdr      []byte `json:"envelope_xdr,omitempty"`
	ResultXdr        []byte `json:"result_xdr,omitempty"`
	ResultMetaXdr    []byte `json:"result_meta_xdr,omitempty"`
	SourceAddress    string `json:"source_address,omitempty"`
}

type ScAddress struct {
	AccountId  *string `json:"account_id,omitempty"`
	ContractId *string `json:"contract_id,omitempty"`
}

type Contract struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	AccountId  string `json:"account_id,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	Ledger     uint32 `json:"ledger,omitempty"`
	EntryType  string `json:"entry_type,omitempty"`
	KeyXdr     []byte `json:"key_xdr,omitempty"`
	ValueXdr   []byte `json:"value_xdr,omitempty"`
	Durability int32  `json:"durability,omitempty"`
	IsNewest   bool   `json:"is_newest,omitempty"`
}

type Int128Parts struct {
	Hi int64  `json:"hi,omitempty"`
	Lo uint64 `json:"lo,omitempty"`
}
type WasmContractEvent struct {
	Id           string `json:"id,omitempty"`
	ContractId   string `json:"contract_id,omitempty"`
	TxHash       string `json:"tx_hash,omitempty"`
	EventBodyXdr []byte `json:"event_body_xdr,omitempty"`
}

type StellarAssetContractEvent interface {
	GetType() string
	Parse(topics xdr.ScVec, value xdr.ScVal) error
}

var (
	ErrNotBalanceChangeEvent = errors.New("event doesn't represent a balance change")
	ErrNotTransferEvent      = errors.New("this is not transfer event")
	ErrNotMintEvent          = errors.New("this is not mint event")
	ErrNotClawbackEvent      = errors.New("this is not clawback event")
	ErrNotBurnEvent          = errors.New("this is not burn event")
)

type AssetContractTransferEvent struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	From       string `json:"from,omitempty"`
	To         string `json:"to,omitempty"`
	AmountHi   int64  `json:"amount_hi,omitempty"`
	AmountLo   uint64 `json:"amount_lo,omitempty"`
}

func (AssetContractTransferEvent) GetType() string {
	return "transfer"
}

func (a *AssetContractTransferEvent) Parse(topics xdr.ScVec, value xdr.ScVal) error {
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
	a.From, a.To, a.AmountHi, a.AmountLo, err = parseBalanceChangeEvent(topics, value)
	if err != nil {
		return ErrNotTransferEvent
	}
	return nil
}

type AssetContractMintEvent struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	Admin      string `json:"admin,omitempty"`
	To         string `json:"to,omitempty"`
	AmountHi   int64  `json:"amount_hi,omitempty"`
	AmountLo   uint64 `json:"amount_lo,omitempty"`
}

func (AssetContractMintEvent) GetType() string {
	return "mint"
}

func (a *AssetContractMintEvent) Parse(topics xdr.ScVec, value xdr.ScVal) error {
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
	a.Admin, a.To, a.AmountHi, a.AmountLo, err = parseBalanceChangeEvent(topics, value)
	if err != nil {
		return ErrNotTransferEvent
	}
	return nil
}

type AssetContractBurnEvent struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	From       string `json:"from,omitempty"`
	AmountHi   int64  `json:"amount_hi,omitempty"`
	AmountLo   uint64 `json:"amount_lo,omitempty"`
}

func (AssetContractBurnEvent) GetType() string {
	return "burn"
}

func (event *AssetContractBurnEvent) Parse(topics xdr.ScVec, value xdr.ScVal) error {
	//
	// The burn event format is:
	//
	// 	"burn"  	Symbol
	//  <from> 		Address
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
	val := XdrInt128PartsConvert(amount)

	event.AmountHi = val.Hi
	event.AmountLo = val.Lo

	return nil
}

type AssetContractClawbackEvent struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	Admin      string `json:"admin,omitempty"`
	From       string `json:"from,omitempty"`
	AmountHi   int64  `json:"amount_hi,omitempty"`
	AmountLo   uint64 `json:"amount_lo,omitempty"`
}

func (AssetContractClawbackEvent) GetType() string {
	return "clawback"
}

func (a *AssetContractClawbackEvent) Parse(topics xdr.ScVec, value xdr.ScVal) error {
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
	a.Admin, a.From, a.AmountHi, a.AmountLo, err = parseBalanceChangeEvent(topics, value)
	if err != nil {
		return ErrNotTransferEvent
	}
	return nil
}

// parseBalanceChangeEvent is a generalization of a subset of the Stellar Asset
// Contract events. Transfer, mint, clawback, and burn events all have two
// addresses and an amount involved. The addresses represent different things in
// different event types (e.g. "from" or "admin"), but the parsing is identical.
// This helper extracts all three parts or returns a generic error if it can't.
func parseBalanceChangeEvent(topics xdr.ScVec, value xdr.ScVal) (
	first string,
	second string,
	amountHi int64,
	amountlo uint64,
	err error,
) {
	err = ErrNotBalanceChangeEvent
	if len(topics) != 4 {
		err = fmt.Errorf("")
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

	amount := XdrInt128PartsConvert(xdrAmount)

	return first, second, amount.Hi, amount.Lo, nil
}

func XdrInt128PartsConvert(in xdr.Int128Parts) Int128Parts {
	out := Int128Parts{
		Hi: int64(in.Hi),
		Lo: uint64(in.Lo),
	}

	return out
}
