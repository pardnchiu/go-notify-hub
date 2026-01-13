package handler

import (
	"encoding/json"
	"log/slog"
	"maps"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"

	"goNotify/internal/utils"
)

// POST: /discord/add
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
		if !vaildDiscordWebhook.MatchString(webhook) {
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

	discordChannelsMu.Lock()
	if discordChannels == nil {
		if data, err := os.ReadFile(path); err == nil {
			var m map[string]string
			if err := json.Unmarshal(data, &m); err == nil {
				discordChannels = m
			} else {
				discordChannels = make(map[string]string)
			}
		} else {
			if os.IsNotExist(err) {
				discordChannels = make(map[string]string)
			} else {
				discordChannelsMu.Unlock()
				slog.Error("Failed to read discord_channel.json", "path", path, "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load channel configuration"})
				return
			}
		}
	}

	for _, data := range req.Datas {
		name := strings.TrimSpace(data.Name)
		webhook := strings.TrimSpace(data.Webhook)

		discordChannels[name] = webhook
	}

	newContent := make(map[string]string, len(discordChannels))
	maps.Copy(newContent, discordChannels)
	discordChannelsMu.Unlock()

	if err := utils.WriteJSON(path, newContent); err != nil {
		slog.Error("Failed to write discord_channel.json", "path", path, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save channel configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "channels added successfully"})
}
