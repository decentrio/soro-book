package models

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
	EnvelopeXdr      string `json:"envelope_xdr,omitempty"`
	ResultXdr        string `json:"result_xdr,omitempty"`
	ResultMetaXdr    string `json:"result_meta_xdr,omitempty"`
	SourceAddress    string `json:"source_address,omitempty"`
}

type ContractEvent struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	LedgerSeq  uint32 `json:"ledger_seq,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	EventType  string `json:"event_type,omitempty"`
	Value      string `json:"value,omitempty"`
}

type ContractEventJSON struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	LedgerSeq  uint32 `json:"ledger_seq,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	EventType  string `json:"event_type,omitempty"`
	Data       []byte `json:"data,omitempty"`
}

type Topics struct {
	EventId    string `json:"event_id,omitempty"`
	TopicXdr   string `json:"topic_xdr,omitempty"`
	TopicIndex int32  `json:"topic_index,omitempty"`
}
