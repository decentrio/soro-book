package handlers

import (
	"github.com/decentrio/soro-book/database/models"
)

func (h *DBHandler) CreateEvent(data *models.Event) (string, error) {
	if err := h.db.Create(data).Error; err != nil {
		return "", err
	}

	return data.ID, nil
}
