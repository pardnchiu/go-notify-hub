package handler

import (
	"encoding/json"
	"goNotify/internal/channel"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

func NewDiscordHandler() *DiscordHandler {
	return &DiscordHandler{}
}

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

// POST: /discord/add/
// BODY: { datas: [{ "name": "name", "webhook": "url"}] }
func (h *DiscordHandler) Add(c *gin.Context) {
	var req struct {
		Datas []struct {
			Name    string `json:"name"`
			Webhook string `json:"webhook"`
		} `json:"datas"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if len(req.Datas) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no channel data provided"})
		return
	}

	var invalidChannelNames []string
	var invalodWebhookURLs []string
	for _, data := range req.Datas {
		name := strings.TrimSpace(data.Name)
		webhook := strings.TrimSpace(data.Webhook)

		if !validChannelName.MatchString(name) {
			slog.Error("Invalid channel name format", "channelName", name)
			invalidChannelNames = append(invalidChannelNames, name)
		}
		if !vaildWebhookURL.MatchString(webhook) {
			slog.Error("Invalid webhook URL format", "webhook", webhook)
			invalodWebhookURLs = append(invalodWebhookURLs, webhook)
		}
	}

	if len(invalidChannelNames) > 0 || len(invalodWebhookURLs) > 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":                 "invalid channel names or webhook URLs",
			"invalid_channel_names": invalidChannelNames,
			"invalid_webhook_urls":  invalodWebhookURLs,
		})
		return
	}

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
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load channel configuration"})
				return
			}
		}
	}

	for _, data := range req.Datas {
		name := strings.TrimSpace(data.Name)
		webhook := strings.TrimSpace(data.Webhook)

		channels[name] = webhook
	}

	toWrite := make(map[string]string, len(channels))
	for k, v := range channels {
		toWrite[k] = v
	}
	channelsMu.Unlock()

	if err := writeJSON(path, toWrite); err != nil {
		slog.Error("Failed to write discord_channel.json", "path", path, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save channel configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "channels added successfully"})
}

// DELETE: /discord/delete/:channelName
func (h *DiscordHandler) Delete(c *gin.Context) {
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

	channelsMu.Lock()
	defer channelsMu.Unlock()

	if channels == nil {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
		return
	}

	delete(channels, channelName)
	toWrite := make(map[string]string, len(channels))
	for k, v := range channels {
		toWrite[k] = v
	}

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

	if err := writeJSON(path, toWrite); err != nil {
		slog.Error("Failed to write discord_channel.json", "path", path, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save channel configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "channel deleted successfully"})
}

func writeJSON(path string, data map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, dataBytes, 0644); err != nil {
		return err
	}
	return nil
}
