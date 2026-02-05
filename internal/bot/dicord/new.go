package discord

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	session  *discordgo.Session
	commands []*discordgo.ApplicationCommand
}

func New() (*Bot, error) {
	token := os.Getenv("DISCORD_TOKEN")

	if token == "" {
		return nil, nil
	}

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("create discord session: %w", err)
	}

	bot := &Bot{session: session}

	session.AddHandler(bot.handleMessageCreate)
	session.AddHandler(bot.handleInteraction)
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentDirectMessages | discordgo.IntentMessageContent

	if err := session.Open(); err != nil {
		return nil, fmt.Errorf("open websocket connection: %w", err)
	}

	slog.Info("bot is running", slog.String("user", session.State.User.Username))

	return bot, nil
}

func (b *Bot) Close() error {
	slog.Info("shutting down")
	if b.session == nil {
		return nil
	}
	return b.session.Close()
}
