package aggregation

import (
	"fmt"
	"math"

	"github.com/decentrio/soro-book/database/models"
	"github.com/google/uuid"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/xdr"
)

type Tickers struct {
	TickerId       string // PHO_USDC
	BaseCurrency   string // PHO
	TargetCurrency string // USDC
	PoolId         string // "CAZ6W4WHVGQBGURYTUOLCUOOHW6VQGAAPSPCD72VEDZMBBPY7H43AYEC"
	LastPrice      string // Last price trade
	BaseVolume     string // base currency trade volume (24h)
	TargetVolume   string // target currency trade volume (24h)
	LiquidityInUsd string // liquidity in usd
	UpdatedLedger  uint32 // updated ledger
}

type HistoricalTrades struct {
	TradeId        string // A unique ID associated with the trade for the currency pair transaction
	Price          string // Transaction price of base asset in target currency
	BaseVolume     string // volume trade of base currency (float)
	TargetVolume   string // volume trade of target currency (float)
	TradeTimestamp uint64 // time stamp of trade
	TradeType      string // buy/sell
}

func (as *Aggregation) handleReceiveNewContractDataEntry(e models.ContractsData) {
	_, err := as.db.CreateContractEntry(&e)
	if err != nil {
		as.Logger.Error(fmt.Sprintf("Error create contract data entry ledger %d tx %s: %s", e.Ledger, e.TxHash, err.Error()))
	}
}

func (tw TransactionWrapper) GetModelsContractDataEntry() []models.ContractsData {
	v3 := tw.Tx.UnsafeMeta.V3
	if v3 == nil {
		return nil
	}

	var entries []models.ContractsData
	for _, op := range v3.Operations {
		for _, change := range op.Changes {
			entry, entryType, found := ContractDataEntry(change)
			// continue with "state" because we don't want to store this entry
			if entryType == "state" {
				continue
			}
			if found {
				keyBz, _ := entry.Key.MarshalBinary()

				var valBz []byte
				if entryType != "removed" {
					valBz, _ = entry.Val.MarshalBinary()
				}
				var contractId string
				var err error
				if entry.Contract.ContractId != nil {
					contractId, err = strkey.Encode(strkey.VersionByteContract, entry.Contract.ContractId[:])
					if err != nil {
						continue
					}
				}

				var accountId string
				if entry.Contract.AccountId != nil {
					accountId, err = entry.Contract.AccountId.GetAddress()
					if err != nil {
						continue
					}
				}

				entry := models.ContractsData{
					Id:            uuid.New().String(),
					ContractId:    contractId,
					AccountId:     accountId,
					TxHash:        tw.GetTransactionHash(),
					Ledger:        tw.GetLedgerSequence(),
					EntryType:     entryType,
					KeyXdr:        keyBz,
					ValueXdr:      valBz,
					Durability:    int32(entry.Durability),
					IsNewest:      true,
					UpdatedLedger: uint32(math.MaxInt32),
				}
				entries = append(entries, entry)
			}
		}
	}

	return entries
}

func ContractDataEntry(c xdr.LedgerEntryChange) (xdr.ContractDataEntry, string, bool) {
	var result xdr.ContractDataEntry

	switch c.Type {
	case xdr.LedgerEntryChangeTypeLedgerEntryCreated:
		created := *c.Created
		if created.Data.ContractData != nil {
			result = *created.Data.ContractData
			return result, "created", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryUpdated:
		updated := *c.Updated
		if updated.Data.ContractData != nil {
			result = *updated.Data.ContractData
			return result, "updated", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryRemoved:
		ledgerKey := c.Removed
		if ledgerKey.ContractData != nil {
			result.Contract = ledgerKey.ContractData.Contract
			result.Key = ledgerKey.ContractData.Key
			result.Durability = ledgerKey.ContractData.Durability
			return result, "removed", true
		}
	case xdr.LedgerEntryChangeTypeLedgerEntryState:
		state := *c.State
		if state.Data.ContractData != nil {
			result = *state.Data.ContractData
			return result, "state", true
		}

	}
	return result, "", false
}
