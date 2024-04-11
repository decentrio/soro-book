package handlers

import (
	"fmt"

	"github.com/decentrio/soro-book/database/models"
)

func (h *DBHandler) CreateLedger(data *models.Ledger) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.Hash, nil
}

func (h *DBHandler) CreateTransaction(data *models.Transaction) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.Hash, nil
}

func (h *DBHandler) CreateWasmContractEvent(data *models.WasmContractEvent) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.Id, nil
}

func (h *DBHandler) CreateAssetContractTransferEvent(data *models.AssetContractTransferEvent) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.Id, nil
}

func (h *DBHandler) CreateAssetContractMintEvent(data *models.AssetContractMintEvent) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.Id, nil
}

func (h *DBHandler) CreateAssetContractBurnEvent(data *models.AssetContractBurnEvent) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.Id, nil
}

func (h *DBHandler) CreateAssetContractClawbackEvent(data *models.AssetContractClawbackEvent) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.Id, nil
}

func (h *DBHandler) CreateContractEntry(data *models.Contract) (string, error) {
	switch data.EntryType {
	case "updated":
		var oldData models.Contract
		if err := h.db.Table("contracts").
			Where("contract_id = ?", data.ContractId).
			Where("is_newest = ?", true).
			Where("key_xdr = ?", data.KeyXdr).
			First(&oldData).Error; err == nil {
			fmt.Println("CreateContractEntry Updated")
			oldData.IsNewest = false
			if err := h.db.Table("contracts").Save(oldData).Error; err != nil {
				return "ERROR: update old contract data entry", err
			}
		}

		break
	case "removed":
		var oldData models.Contract
		if err := h.db.Table("contracts").
			Where("contract_id = ?", data.ContractId).
			Where("is_newest = ?", true).
			Where("key_xdr = ?", data.KeyXdr).
			First(&oldData).Error; err == nil {
			oldData.IsNewest = false
			if err := h.db.Table("contracts").Save(oldData).Error; err != nil {
				return "ERROR: update old contract data entry", err
			}
		}

		break
	}

	if err := h.db.Create(data).Error; err != nil {
		return "ERROR: create contract data entry", err
	}

	return fmt.Sprintf("%s: %s-%s", data.EntryType, data.ContractId, string(data.KeyXdr)), nil
}
