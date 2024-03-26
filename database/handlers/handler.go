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

func (h *DBHandler) CreateEvent(data *models.Event) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.Id, nil
}

func (h *DBHandler) CreateTopics(data *models.Topics) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return fmt.Sprintf("%s-%d", data.EventId, data.TopicIndex), nil
}
