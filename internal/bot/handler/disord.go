package handler

import (
	"context"

	"github.com/bwmarrin/discordgo"
)

type DiscordReply struct {
	Session     *discordgo.Session
	Interaction *discordgo.InteractionCreate
	ChannelID   string
	Reference   *discordgo.MessageReference
	IsFirst     bool
}

func ReplyDiscord(ctx context.Context, dr *DiscordReply, reply Reply) error {
	var embeds []*discordgo.MessageEmbed

	if reply.ImageURL != "" {
		embeds = []*discordgo.MessageEmbed{
			{
				Image: &discordgo.MessageEmbedImage{
					URL: reply.ImageURL,
				},
			},
		}
	}

	if dr.Interaction != nil {
		if dr.IsFirst {
			return dr.Session.InteractionRespond(dr.Interaction.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: reply.Content,
					Embeds:  embeds,
				},
			})
		}
		_, err := dr.Session.FollowupMessageCreate(dr.Interaction.Interaction, true, &discordgo.WebhookParams{
			Content: reply.Content,
			Embeds:  embeds,
		})
		return err
	}

	_, err := dr.Session.ChannelMessageSendComplex(dr.ChannelID, &discordgo.MessageSend{
		Content:   reply.Content,
		Reference: dr.Reference,
		Embeds:    embeds,
	})
	return err
}
