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

func (h *DBHandler) CreateContractEntry(data *models.Contract) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%s", data.ContractId, string(data.KeyXdr)), nil
}
