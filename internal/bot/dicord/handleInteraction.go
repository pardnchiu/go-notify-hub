package discord

import (
	"context"
	"go-notify-hub/internal/bot/handler"
	"log/slog"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func (b *Bot) handleInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	data := i.ApplicationCommandData()

	var userID, username string
	if i.Member != nil {
		userID = i.Member.User.ID
		username = i.Member.User.Username
	} else if i.User != nil {
		userID = i.User.ID
		username = i.User.Username
	}

	slog.Info("slash command received",
		slog.String("command", data.Name),
		slog.String("user_id", userID),
		slog.String("user", username),
	)

	msg := parseInteraction(i)
	replies := handler.Handler(msg)

	ctx := context.Background()
	for idx, reply := range replies {
		dr := &handler.DiscordReply{
			Session:     s,
			Interaction: i,
			IsFirst:     idx == 0,
		}
		handler.ReplyDiscord(ctx, dr, reply)
	}
}

func parseInteraction(i *discordgo.InteractionCreate) *handler.Message {
	data := i.ApplicationCommandData()

	var userID, channelID, guildID string
	if i.Member != nil {
		userID = i.Member.User.ID
		guildID = i.GuildID
	} else if i.User != nil {
		userID = i.User.ID
	}
	channelID = i.ChannelID

	var params []string
	for _, opt := range data.Options {
		params = append(params, opt.StringValue())
	}

	return &handler.Message{
		UserID:    userID,
		ChannelID: channelID,
		GuildID:   guildID,
		Cmd:       "/" + data.Name,
		Params:    params,
		Content:   "/" + data.Name + " " + strings.Join(params, " "),
	}
}
