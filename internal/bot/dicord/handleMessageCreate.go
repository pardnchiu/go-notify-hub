package discord

import (
	"context"
	"fmt"
	"go-notify-hub/internal/bot/handler"
	"log/slog"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	slog.Info("message received",
		slog.String("guild_id", m.GuildID),
		slog.String("channel_id", m.ChannelID),
		slog.String("author_id", m.Author.ID),
		slog.String("author", m.Author.Username),
		slog.String("content", m.Content),
	)

	msg, err := parseMessage(m)
	if err != nil {
		return
	}

	replies := handler.Handler(msg)

	ctx := context.Background()
	for _, reply := range replies {
		dr := &handler.DiscordReply{
			Session:   s,
			ChannelID: m.ChannelID,
			Reference: m.Reference(),
		}
		handler.ReplyDiscord(ctx, dr, reply)
	}
}

func parseMessage(m *discordgo.MessageCreate) (*handler.Message, error) {
	content := strings.TrimSpace(m.Content)
	regex := regexp.MustCompile(`^/\w+(\s+\S+)*$`)
	if !regex.MatchString(content) {
		return nil, fmt.Errorf("not a command")
	}

	fields := strings.Fields(content)
	newMsg := &handler.Message{
		UserID:    m.Author.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
		Cmd:       fields[0],
		Params:    fields[1:],
		Content:   content,
	}

	return newMsg, nil
}
