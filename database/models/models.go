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
	EnvelopeXdr      []byte `json:"envelope_xdr,omitempty"`
	ResultXdr        []byte `json:"result_xdr,omitempty"`
	ResultMetaXdr    []byte `json:"result_meta_xdr,omitempty"`
	SourceAddress    string `json:"source_address,omitempty"`
}
type Event struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	TxIndex    uint32 `json:"tx_index,omitempty"`
	EventXdr   []byte `json:"event_xdr,omitempty"`
}

type Topics struct {
	EventId  string `json:"event_id,omitempty"`
	TopicXdr []byte `json:"topic_xdr,omitempty"`
	TopicIdx int32  `json:"topic_idx,omitempty"`
}
