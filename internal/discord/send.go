package discord

import (
	"go-notification-bot/internal/channel"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

// POST: /discord/send/:channelName
func (h *DiscordHandler) Send(c *gin.Context) {
	channelName := c.Param("channelName")
	if channelName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "channel name is required"})
		return
	}

	if !validChannelName.MatchString(channelName) {
		slog.Error("Invalid channel name format", "channelName", channelName)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid channel name format"})
		return
	}

	discordChannelsMu.RLock()
	cacheChannels := discordChannels
	discordChannelsMu.RUnlock()
	if cacheChannels == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "please insert this channel first"})
		return
	}

	webhook, ok := cacheChannels[channelName]
	if !ok || webhook == "" {
		slog.Error("Channel does not exist or has empty webhook", "channelName", channelName)
		c.JSON(http.StatusBadRequest, gin.H{"error": "this channel does not exist"})
		return
	}

	var req channel.DiscordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}
	req.WebhookURL = webhook

	if req.Title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title is required"})
		return
	}

	if req.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "description is required"})
		return
	}

	if err := channel.SendToDiscord(req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "notification sent successfully"})
}
