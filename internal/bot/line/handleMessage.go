package Linebot

import (
	"context"
	"fmt"
	"go-notify-hub/internal/bot/handler"
	"regexp"
	"strings"

	"github.com/line/line-bot-sdk-go/v8/linebot"
)

func (h *LinebotHandler) handleMessage(ctx context.Context, event *linebot.Event, bot *linebot.Client) {
	if event.Message == nil {
		return
	}

	msg, ok := event.Message.(*linebot.TextMessage)
	if !ok {
		return
	}

	userID := event.Source.UserID
	if userID == "" {
		return
	}

	parse, err := parseMessage(userID, msg.Text)
	if err != nil {
		return
	}

	messages := handler.Handler(parse)

	for _, e := range messages {
		handler.ReplyLine(ctx, event, bot, e)
	}
}

func parseMessage(userID, msg string) (*handler.Message, error) {
	content := strings.TrimSpace(msg)
	regex := regexp.MustCompile(`^/\w+(\s+\S+)*$`)
	if !regex.MatchString(content) {
		return nil, fmt.Errorf("invalid command format")
	}

	fields := strings.Fields(content)
	newMsg := &handler.Message{
		UserID:  userID,
		Cmd:     fields[0],
		Params:  fields[1:],
		Content: msg,
	}

	return newMsg, nil
}
