package main

import (
	"goNotify/internal/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	discordHandler, err := handler.NewDiscordHandler()
	if err != nil {
		log.Fatal("Failed to create Discord handler:", err)
	}

	slackHandler, err := handler.NewSlackHandler()
	if err != nil {
		log.Fatal("Failed to create Slack handler:", err)
	}

	r.GET("/discord/list", discordHandler.List)
	r.POST("/discord/:channelName", discordHandler.Send)
	r.POST("/discord/add", discordHandler.Add)
	r.DELETE("/discord/:channelName", discordHandler.Delete)

	r.GET("/slack/list", slackHandler.List)

	log.Println("start on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
