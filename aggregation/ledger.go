package aggregation

import (
	"github.com/decentrio/soro-book/database/models"
	"github.com/stellar/go/xdr"
)

func getLedgerFromCloseMeta(ledgerCloseMeta xdr.LedgerCloseMeta) models.Ledger {
	return models.Ledger{
		Hash:     ledgerCloseMeta.LedgerHash().HexString(),
		PrevHash: ledgerCloseMeta.PreviousLedgerHash().HexString(),
		Sequence: ledgerCloseMeta.LedgerSequence(),
	}
}
