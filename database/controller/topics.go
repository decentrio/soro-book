package controller

import (
	"net/http"

	"github.com/decentrio/soro-book/database/handlers"
	"github.com/decentrio/soro-book/database/models"
	"github.com/gin-gonic/gin"
)

func CreateTopics(h *handlers.DBHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var dataItem models.Topics

		if err := c.BindJSON(&dataItem); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		topicIndex, err := h.CreateTopics(&dataItem)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"data": topicIndex})
	}
}
