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
type WasmContractEvent struct {
	Id         string `json:"id,omitempty"`
	ContractId string `json:"contract_id,omitempty"`
	TxHash     string `json:"tx_hash,omitempty"`
	EventXdr   []byte `json:"event_xdr,omitempty"`
}

type Int128Parts struct {
	Hi int64  `json:"hi,omitempty"`
	Lo uint64 `json:"lo,omitempty"`
}

type AssetContractTransferEvent struct {
	Id         string      `json:"id,omitempty"`
	ContractId string      `json:"contract_id,omitempty"`
	TxHash     string      `json:"tx_hash,omitempty"`
	From       string      `json:"from,omitempty"`
	To         string      `json:"to,omitempty"`
	Amount     Int128Parts `json:"amount,omitempty"`
}

type AssetContractMintEvent struct {
	Id         string      `json:"id,omitempty"`
	ContractId string      `json:"contract_id,omitempty"`
	TxHash     string      `json:"tx_hash,omitempty"`
	Admin      string      `json:"admin,omitempty"`
	To         string      `json:"to,omitempty"`
	Amount     Int128Parts `json:"amount,omitempty"`
}

type AssetContractBurnEvent struct {
	Id         string      `json:"id,omitempty"`
	ContractId string      `json:"contract_id,omitempty"`
	TxHash     string      `json:"tx_hash,omitempty"`
	From       string      `json:"from,omitempty"`
	Amount     Int128Parts `json:"amount,omitempty"`
}

type AssetContractClawbackEvent struct {
	Id         string      `json:"id,omitempty"`
	ContractId string      `json:"contract_id,omitempty"`
	TxHash     string      `json:"tx_hash,omitempty"`
	Admin      string      `json:"admin,omitempty"`
	From       string      `json:"from,omitempty"`
	Amount     Int128Parts `json:"amount,omitempty"`
}

type Contract struct {
	ContractId string `json:"contract_id,omitempty"`
	AccountId  string `json:"account_id,omitempty"`
	Ledger     uint32 `json:"ledger,omitempty"`
	KeyXdr     []byte `json:"key_xdr,omitempty"`
	ValueXdr   []byte `json:"value_xdr,omitempty"`
	Durability int32  `json:"durability,omitempty"`
}
