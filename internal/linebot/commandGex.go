package Linebot

import (
	"context"
	"fmt"
	"go-notification-bot/internal/database"
	"log/slog"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot"
)

func (h *LinebotHandler) commandGex(ctx context.Context, msg *Message, event *linebot.Event, bot *linebot.Client) {
	const fn = "LinebotHandler/commandGex"
	ticker := strings.ToUpper(msg.Params[0])
	if !vaildTicker.MatchString(ticker) {
		text := "Usage: /gex $<ticker>"
		_, err := bot.ReplyMessage(
			event.ReplyToken,
			linebot.NewTextMessage(text),
		).WithContext(ctx).Do()
		if err != nil {
			slog.Error(fn+": failed to reply format error", "error", err)
		}
		return
	}
	ticker = strings.TrimPrefix(ticker, "$")

	result, err := database.SelectTicker(ticker)
	if err != nil {
		slog.Error(fn+": failed to get ticker data", "ticker", ticker, "error", err)
		return
	}

	_, err = bot.ReplyMessage(
		event.ReplyToken,
		linebot.NewImageMessage(
			result.ImageURL+"?width=2048",
			result.ImageURL,
		),
		linebot.NewTextMessage(fmt.Sprintf(`親愛的投資人您好，今天是 %s，
盤前觀察 $%s
%s

%s

%s`,
			time.Now().Format("2006-01-02"),
			ticker,
			result.RationaleCompact,
			result.ExpiryNearTerm,
			result.ExpiryFarTerm,
		)),
	).WithContext(ctx).Do()
	if err != nil {
		slog.Error(fn+": failed to reply ticker result", "userID", msg.UserID, "error", err)
	}
}
