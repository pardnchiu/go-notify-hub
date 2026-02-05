package main

import (
	discordbot "go-notify-hub/internal/bot/dicord"
	Linebot "go-notify-hub/internal/bot/line"
	"go-notify-hub/internal/channel/discord"
	"go-notify-hub/internal/channel/slack"
	"go-notify-hub/internal/database"
	"go-notify-hub/internal/email"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		slog.Warn("failed to load .env",
			slog.String("error", err.Error()))
	}
}

func main() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		if _, err := os.Stat("/.dockerenv"); err == nil {
			dbPath = "/data/database.db"
		} else {
			home, _ := os.UserHomeDir()
			dbPath = filepath.Join(home, ".go-notify-hub", "database.db")
			dir := filepath.Dir(dbPath)
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalf("Failed to create directory %s: %v", dir, err)
			}
		}
	}

	db, err := database.NewSQLite(dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	bot, err := discordbot.New()
	if bot != nil {
		defer bot.Close()
	}

	r := gin.Default()

	discordHandler, err := discord.New()
	if err != nil {
		log.Fatal("Failed to create Discord handler:", err)
	}

	r.GET("/discord/list", discordHandler.List)
	r.POST("/discord/:channelName", discordHandler.Send)
	r.POST("/discord/add", discordHandler.Add)
	r.DELETE("/discord/:channelName", discordHandler.Delete)

	slackHandler, err := slack.New()
	if err != nil {
		log.Fatal("Failed to create Slack handler:", err)
	}

	r.GET("/slack/list", slackHandler.List)
	r.POST("/slack/:channelName", slackHandler.Send)
	r.POST("/slack/add", slackHandler.Add)
	r.DELETE("/slack/:channelName", slackHandler.Delete)

	secret := os.Getenv("LINEBOT_SECRET")
	token := os.Getenv("LINEBOT_TOKEN")
	if secret != "" && token != "" {
		linebotHandler, err := Linebot.New()
		if err != nil {
			log.Fatal("Failed to create linebot handler:", err)
		}
		r.POST("/linebot/webhook", linebotHandler.Webhook)
		r.POST("/linebot/send/all", linebotHandler.Send)
	}

	emailHandler, err := email.New()
	if err != nil {
		log.Fatal("Failed to create email handler:", err)
	}
	r.POST("/email/send", emailHandler.Send)
	r.POST("/email/send/bulk", emailHandler.SendBulk)

	r.NoRoute(func(c *gin.Context) {
		select {}
	})

	log.Println("start on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
