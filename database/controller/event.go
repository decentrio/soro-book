package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/decentrio/soro-book/database/handlers"
	"github.com/decentrio/soro-book/database/models"
)

func CreateEvent(h *handlers.DBHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataItem models.WasmContractEvent

		if err := c.BindJSON(&dataItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		eventID, err := h.CreateWasmContractEvent(&dataItem)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": eventID})
	}
}

func HelloEvent(_ *handlers.DBHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"hello": "ok"})
	}
}
