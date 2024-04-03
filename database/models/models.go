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

type Contract struct {
	ContractId          string `json:"contract_id,omitempty"`
	AccountId           string `json:"account_id,omitempty"`
	ExpirationLedgerSeq uint32 `json:"expiration_ledger_seq,omitempty"`
	KeyXdr              []byte `json:"key_xdr,omitempty"`
	ValueXdr            []byte `json:"value_xdr,omitempty"`
	Durability          int32  `json:"durability,omitempty"`
}
