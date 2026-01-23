package Linebot

import (
	"context"
	"fmt"
	"log/slog"
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

	parseMsg, err := parseMessage(userID, msg.Text)
	if err != nil {
		return
	}

	if err := verifyMessage(*parseMsg); err != nil {
		_, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(err.Error()),
		).WithContext(ctx).Do()
		if err != nil {
			slog.Error("LinebotHandler/handleMessage: failed to reply message", "error", err)
		}
		return
	}
}

type Message struct {
	UserID  string
	Cmd     string
	Params  []string
	Message string
}

func parseMessage(userID, msg string) (*Message, error) {
	if !vaildCommand.MatchString(msg) {
		return nil, fmt.Errorf("invalid command format")
	}

	fields := strings.Fields(msg)
	newMsg := &Message{
		UserID:  userID,
		Cmd:     fields[0],
		Params:  fields[1:],
		Message: msg,
	}

	return newMsg, nil
}

func verifyMessage(msg Message) error {
	if msg.Cmd == "/gex" && len(msg.Params) < 1 {
		return fmt.Errorf("Usage: /gex $<ticker>")
	}
	return nil
}
