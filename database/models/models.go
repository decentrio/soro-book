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

type Contract struct {
	ContractId string `json:"contract_id,omitempty"`
	AccountId  string `json:"account_id,omitempty"`
	Ledger     uint32 `json:"ledger,omitempty"`
	KeyXdr     []byte `json:"key_xdr,omitempty"`
	ValueXdr   []byte `json:"value_xdr,omitempty"`
	Durability int32  `json:"durability,omitempty"`
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
}

var (
	ErrNotBalanceChangeEvent = errors.New("event doesn't represent a balance change")
	ErrNotTransferEvent      = errors.New("this is not transfer event")
	ErrNotMintEvent          = errors.New("this is not mint event")
	ErrNotClawbackEvent      = errors.New("this is not clawback event")
	ErrNotBurnEvent          = errors.New("this is not burn event")
)

type AssetContractTransferEvent struct {
	Id           string `json:"id,omitempty"`
	ContractId   string `json:"contract_id,omitempty"`
	TxHash       string `json:"tx_hash,omitempty"`
	EventBodyXdr []byte `json:"event_body_xdr,omitempty"`
}

func (AssetContractTransferEvent) GetType() string {
	return "transfer"
}

type AssetContractMintEvent struct {
	Id           string `json:"id,omitempty"`
	ContractId   string `json:"contract_id,omitempty"`
	TxHash       string `json:"tx_hash,omitempty"`
	EventBodyXdr []byte `json:"event_body_xdr,omitempty"`
}

func (AssetContractMintEvent) GetType() string {
	return "mint"
}

type AssetContractBurnEvent struct {
	Id           string `json:"id,omitempty"`
	ContractId   string `json:"contract_id,omitempty"`
	TxHash       string `json:"tx_hash,omitempty"`
	EventBodyXdr []byte `json:"event_body_xdr,omitempty"`
}

func (AssetContractBurnEvent) GetType() string {
	return "burn"
}

type AssetContractClawbackEvent struct {
	Id           string `json:"id,omitempty"`
	ContractId   string `json:"contract_id,omitempty"`
	TxHash       string `json:"tx_hash,omitempty"`
	EventBodyXdr []byte `json:"event_body_xdr,omitempty"`
}

func (AssetContractClawbackEvent) GetType() string {
	return "clawback"
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
