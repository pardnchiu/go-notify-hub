package Linebot

import (
	"context"
	"go-notify-hub/internal/database"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
)

type LinebotMessage struct {
	Text         string `json:"text,omitempty"`
	Image        string `json:"image,omitempty"`
	ImagePreview string `json:"image_preview,omitempty"`
}

// * POST: /linebot/send/all
func (h *LinebotHandler) Send(c *gin.Context) {
	const fn = "LinebotHandler/LinebotSend"
	var (
		req LinebotMessage
	)
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error(fn+"[0]", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to parse request payload"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	userIDs, err := database.SelectUserLinebot(ctx)
	if err != nil {
		slog.Error(fn+"[1]", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user IDs"})
		return
	}

	// * Line boardcast max is 500 messages per request
	if len(userIDs) <= 500 {
		err := send(ctx, userIDs, req.Text, req.Image, req.ImagePreview)
		if err != nil {
			slog.Error(fn+"[2]", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
		}
		return
	}

	for i := 0; i < len(userIDs); i += 500 {
		end := min(i+500, len(userIDs))
		err := send(ctx, userIDs[i:end], req.Text, req.Image, req.ImagePreview)
		if err != nil {
			slog.Error(fn+"[3]", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to send message"})
			break
		}
	}
}

func send(ctx context.Context, userIDs []string, text, image, imagePreview string) error {
	if image == "" {
		_, err := Linebot.Multicast(userIDs,
			linebot.NewTextMessage(text),
		).WithContext(ctx).Do()
		return err
	}

	if imagePreview == "" {
		imagePreview = image
	}
	_, err := Linebot.Multicast(userIDs,
		linebot.NewImageMessage(
			image,
			imagePreview,
		),
		linebot.NewTextMessage(text),
	).WithContext(ctx).Do()
	return err
}
