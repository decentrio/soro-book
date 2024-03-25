package models

type Event struct {
	Type   string `json:"type"`
	Ledger int32  `json:"ledger"`
	ID     string `json:"id"`
}

type ContractEvent struct {
	Id         string   `json:"admin,omitempty"`
	ContractId string   `json:"contract_id,omitempty"`
	LedgerSeq  uint32   `json:"ledger_seq,omitempty"`
	TxHash     string   `json:"tx_hash,omitempty"`
	EventType  string   `json:"type,omitempty"`
	Topics     []string `json:"topics,omitempty"`
	Value      string   `json:"value,omitempty"`
}

type ContractEventJSON struct {
	Id         string `json:"admin,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	LedgerSeq  uint32 `json:"ledger_seq,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	EventType  string `json:"type,omitempty"`
	Data       []byte `json:"data,omitempty"`
}
