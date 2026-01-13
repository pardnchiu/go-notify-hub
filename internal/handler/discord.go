package handler

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	channels         map[string]string
	channelsMu       sync.RWMutex
	validChannelName = regexp.MustCompile(`^[0-9A-Za-z@_-]+$`)
	vaildWebhookURL  = regexp.MustCompile(`^https://discord\.com/api/webhooks/\d{17,20}/[A-Za-z0-9_\-]{68}$`)
)

type DiscordHandler struct{}

func NewDiscordHandler() (*DiscordHandler, error) {
	wd, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to get working directory", "error", err)
		return nil, err
	}

	path := filepath.Join(wd, "json", "discord_channel.json")
	path = filepath.ToSlash(path)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	channelsMu.Lock()
	if channels == nil {
		if data, err := os.ReadFile(path); err == nil {
			var m map[string]string
			if err := json.Unmarshal(data, &m); err == nil {
				channels = m
			} else {
				channels = make(map[string]string)
			}
		} else {
			if os.IsNotExist(err) {
				channels = make(map[string]string)
			} else {
				channelsMu.Unlock()
				slog.Error("Failed to read discord_channel.json", "path", path, "error", err)
				return nil, err
			}
		}
	}

	toWrite := make(map[string]string, len(channels))
	for k, v := range channels {
		toWrite[k] = v
	}
	channelsMu.Unlock()

	return &DiscordHandler{}, nil
}

// GET: /discord/list
func (h *DiscordHandler) List(c *gin.Context) {
	channelsMu.RLock()
	defer channelsMu.RUnlock()

	if channels == nil {
		c.JSON(200, gin.H{"channels": map[string]string{}})
		return
	}
	c.JSON(200, gin.H{"channels": channels})
}
