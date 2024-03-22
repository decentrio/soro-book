package models

import (
	"github.com/stellar/go/xdr"
)

type Event struct {
	Type   string `json:"type"`
	Ledger int32  `json:"ledger"`
	ID     string `json:"id"`
}

type EventTx struct {
	Id           string
	ContractId   xdr.Hash
	LedgerNumber uint32
	TxHash       xdr.Hash
	Type         xdr.ContractEventType
	Topics       []byte
	Data         []byte
}
