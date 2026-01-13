package main

import (
	"goNotify/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	discordHandler := handler.NewDiscordHandler()

	r.POST("/discord/:channelName", discordHandler.Send)
	r.POST("/discord/add", discordHandler.Add)

	log.Println("start on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
