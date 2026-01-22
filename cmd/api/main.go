package main

import (
	"go-notification-bot/internal/database"
	"go-notification-bot/internal/discord"
	Linebot "go-notification-bot/internal/linebot"
	"go-notification-bot/internal/slack"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found")
	}

	db, err := database.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := gin.Default()

	token := os.Getenv("DISCORD_TOKEN")
	bot, err := discord.NewBot(token)
	if err != nil {
		slog.Error("failed to create bot", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer bot.Close()

	discordHandler, err := discord.New()
	if err != nil {
		log.Fatal("Failed to create Discord handler:", err)
	}

	slackHandler, err := slack.New()
	if err != nil {
		log.Fatal("Failed to create Slack handler:", err)
	}

	linebotHandler, err := Linebot.New()
	if err != nil {
		log.Fatal("Failed to create linebot handler:", err)
	}

	r.GET("/discord/list", discordHandler.List)
	r.POST("/discord/:channelName", discordHandler.Send)
	r.POST("/discord/add", discordHandler.Add)
	r.DELETE("/discord/:channelName", discordHandler.Delete)

	r.GET("/slack/list", slackHandler.List)
	r.POST("/slack/:channelName", slackHandler.Send)
	r.POST("/slack/add", slackHandler.Add)
	r.DELETE("/slack/:channelName", slackHandler.Delete)

	r.POST("/linebot/webhook", linebotHandler.Webhook)
	r.POST("/linebot/send/all", linebotHandler.Send)

	r.NoRoute(func(c *gin.Context) {
		select {}
	})

	log.Println("start on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
