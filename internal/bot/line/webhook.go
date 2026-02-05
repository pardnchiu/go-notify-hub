package Linebot

import (
	"context"
	"go-notify-hub/internal/database"
	"log"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

var (
	Linebot      *linebot.Client
	linebotMu    sync.Mutex
	vaildCommand = regexp.MustCompile(`^/([A-Za-z0-9]+)`) // detect command syntax
)

type LinebotHandler struct{}

func New() (*LinebotHandler, error) {
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

// * POST: /linebot/webhook
func (h *LinebotHandler) Webhook(c *gin.Context) {
	events, err := Linebot.ParseRequest(c.Request)
	if err != nil {
		slog.Error("LinebotHandler/Webhook: failed to parse request",
			slog.Any("error", err),
		)
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

func (h *LinebotHandler) handleFollow(ctx context.Context, event *linebot.Event) {
	const fn = "LinebotHandler/handleFollow"

	userID := event.Source.UserID
	if userID == "" {
		slog.Error(fn + ": userID is required")
		return
	}

	err := database.DB.InsertUser(ctx, userID)
	if err != nil {
		slog.Error(fn+": failed to insert user",
			slog.String("userID", userID),
			slog.Any("error", err))
		return
	}

	slog.Info("new user added",
		slog.String("userID", userID))
}

func (h *LinebotHandler) handleUnfollow(ctx context.Context, event *linebot.Event) {
	const fn = "LinebotHandler/handleUnfollow"

	userID := event.Source.UserID
	if userID == "" {
		slog.Error(fn + ": userID is required")
		return
	}

	err := database.DB.DeleteUser(ctx, userID)
	if err != nil {
		slog.Error(fn+": failed to delete user",
			slog.String("userID", userID),
			slog.Any("error", err))
		return
	}

	slog.Info("user deleted",
		slog.String("userID", userID))
}
