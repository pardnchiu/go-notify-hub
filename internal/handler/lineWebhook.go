package handler

import (
	"context"
	"fmt"
	"goNotify/internal/database"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

var (
	Linebot      *linebot.Client
	linebotMu    sync.Mutex
	vaildCommand = regexp.MustCompile(`^/([A-Za-z0-9]+)`) // detect command syntax
	vaildTicker  = regexp.MustCompile(`\$[A-Za-z]{3,5}`)  // check stock ticker
)

type LinebotHandler struct{}

func NewLineHandler() (*LinebotHandler, error) {
	linebotMu.Lock()
	defer linebotMu.Unlock()

	if Linebot == nil {
		secret := os.Getenv("LINEBOT_SECRET")
		token := os.Getenv("LINEBOT_TOKEN")

		if secret == "" || token == "" {
			log.Fatal("LINEBOT_SECRET and LINEBOT_TOKEN must be set")
		}

		bot, err := linebot.New(secret, token)
		if err != nil {
			log.Fatal(err)
		}
		Linebot = bot
	}

	return &LinebotHandler{}, nil
}

func (h *LinebotHandler) Webhook(c *gin.Context) {
	events, err := Linebot.ParseRequest(c.Request)
	if err != nil {
		slog.Error("LinebotHandler/Webhook: failed to parse request", "error", err)
		c.Status(http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeMessage:
			h.handleMessage(ctx, event, Linebot)
		case linebot.EventTypeFollow:
			h.handleFollow(ctx, event)
		case linebot.EventTypeUnfollow:
			h.handleUnfollow(ctx, event)
		}
	}
	c.Status(http.StatusOK)
}

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

	switch parseMsg.Cmd {
	case "/gex":
		h.commandGex(ctx, parseMsg, event, bot)
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

func (h *LinebotHandler) handleFollow(ctx context.Context, event *linebot.Event) {
	const fn = "LinebotHandler/handleFollow"
	userID := event.Source.UserID
	if userID == "" {
		slog.Error(fn + ": userID is required")
		return
	}

	err := database.InsertUser(ctx, userID)
	if err != nil {
		slog.Error(fn+": failed to insert user", "userID", userID, "error", err)
		return
	}

	slog.Info("new user added", "userID", userID)
}

func (h *LinebotHandler) handleUnfollow(ctx context.Context, event *linebot.Event) {
	const fn = "LinebotHandler/handleUnfollow"
	userID := event.Source.UserID
	if userID == "" {
		slog.Error(fn + ": userID is required")
		return
	}

	err := database.DeleteUser(ctx, userID)
	if err != nil {
		slog.Error(fn+": failed to delete user", "userID", userID, "error", err)
		return
	}

	slog.Info("user deleted", "userID", userID)
}
