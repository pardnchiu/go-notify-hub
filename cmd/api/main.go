package main

import (
	"goNotify/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	handler, err := handler.NewDiscordHandler()
	if err != nil {
		log.Fatal("Failed to create Discord handler:", err)
	}

	r.GET("/discord/list", handler.List)
	r.POST("/discord/:channelName", handler.Send)
	r.POST("/discord/add", handler.Add)
	r.DELETE("/discord/:channelName", handler.Delete)

	log.Println("start on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
