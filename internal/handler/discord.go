package handler

import (
	"encoding/json"
	"goNotify/internal/channel"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	channels         map[string]string
	channelsMu       sync.RWMutex
	validChannelName = regexp.MustCompile("^[0-9A-Za-z@_-]+$")
)

type DiscordHandler struct{}

func NewDiscordHandler() *DiscordHandler {
	return &DiscordHandler{}
}

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

	channelsMu.RLock()
	cacheChannels := channels
	channelsMu.RUnlock()

	if cacheChannels == nil {
		wd, err := os.Getwd()
		if err != nil {
			slog.Error("Failed to get working directory", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			return
		}

		path := filepath.Join(wd, "json", "discord_channel.json")
		path = filepath.ToSlash(path)
		if abs, err := filepath.Abs(path); err == nil {
			path = abs
		}

		channelsMu.Lock()
		if channels == nil {
			data, err := os.ReadFile(path)
			if err != nil {
				slog.Error("Failed to read discord_channel.json", "path", path, "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "channel configuration not found"})
				channelsMu.Unlock()
				return
			}
			var tempChannels map[string]string
			if err := json.Unmarshal(data, &tempChannels); err != nil {
				slog.Error("Failed to parse discord_channel.json", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid channel configuration"})
				channelsMu.Unlock()
				return
			}
			channels = tempChannels
		}
		cacheChannels = channels
		channelsMu.Unlock()
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
