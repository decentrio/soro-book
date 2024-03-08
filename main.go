package main

import (
	"github.com/gin-gonic/gin"
	"github.com/decentrio/soro-book/database/handlers"
	 "github.com/decentrio/soro-book/controller"
)

func main() {
	handler := handlers.NewDBHandler()

	router := gin.Default()
	api := router.Group("/api")

	{
		api.POST("/event/create", controller.CreateEvent(handler))
		api.GET("/event/hello", controller.HelloEvent(handler))
	}

	router.Run(":4200")
}
