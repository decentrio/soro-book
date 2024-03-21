package models

type Event struct {
	Type   string `json:"type"`
	Ledger int32  `json:"ledger"`
	ID     string `json:"id"`
}
