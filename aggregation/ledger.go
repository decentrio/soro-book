package aggregation

import (
	"github.com/stellar/go/xdr"
)

type Ledger struct {
	Hash         string `json:"hash,omitempty"`
	PrevHash     string `json:"prev_hash,omitempty"`
	Sequence     uint32 `json:"sequence,omitempty"`
	Transactions uint32 `json:"transactions,omitempty"`
	Operations   uint32 `json:"operations,omitempty"`
}

func getLedgerFromCloseMeta(ledgerCloseMeta xdr.LedgerCloseMeta) *Ledger {
	return &Ledger{
		Hash:     ledgerCloseMeta.LedgerHash().HexString(),
		PrevHash: ledgerCloseMeta.PreviousLedgerHash().HexString(),
		Sequence: ledgerCloseMeta.LedgerSequence(),
	}
}
