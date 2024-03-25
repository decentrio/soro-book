package controller

import (
	"net/http"

	"github.com/decentrio/soro-book/database/handlers"
	"github.com/decentrio/soro-book/database/models"
	"github.com/gin-gonic/gin"
)

func CreateTransaction(h *handlers.DBHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataItem models.Transaction

		if err := c.BindJSON(&dataItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		txHash, err := h.CreateTransaction(&dataItem)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": txHash})
	}
}
