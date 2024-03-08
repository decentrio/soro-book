package models

import (
	"time"
)

type Event struct {
	EventType                string    `json:"type"`
	Ledger                   int32     `json:"ledger"`
	LedgerClosedAt           string    `json:"ledger_closed_at"`
	ContractID               string    `json:"contract_id"`
	ID                       string    `json:"id"`
	PagingToken              string    `json:"paging_token"`
	Topic                    string    `json:"topic"`
	Value                    string    `json:"value"`
	InSuccessfulContractCall bool      `json:"in_successful_contract_call"`
	TransactionHash          string    `json:"tx_hash"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}
