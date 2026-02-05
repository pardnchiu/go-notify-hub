package handler

import (
	"context"

	"github.com/line/line-bot-sdk-go/v8/linebot"
)

func ReplyLine(ctx context.Context, event *linebot.Event, bot *linebot.Client, reply Reply) error {
	var messages []linebot.SendingMessage

	if reply.ImageURL != "" {
		messages = append(messages, linebot.NewImageMessage(
			reply.ImageURL,
			reply.PreviewURL,
		))
	}

	if reply.Content != "" {
		messages = append(messages, linebot.NewTextMessage(reply.Content))
	}

	if len(messages) == 0 {
		return nil
	}

	_, err := bot.ReplyMessage(
		event.ReplyToken,
		messages...,
	).WithContext(ctx).Do()

	return err
}
