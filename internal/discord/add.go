package discord

import (
	"log/slog"
	"maps"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-notify-hub/internal/utils"
)

// * POST: /discord/add
// * BODY: { datas: [{ "name": "name", "webhook": "url"}] }
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

	discordChannelsMu.Lock()
	defer discordChannelsMu.Unlock()

	if discordChannels == nil {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
		return
	}

	for _, data := range req.Datas {
		name := strings.TrimSpace(data.Name)
		webhook := strings.TrimSpace(data.Webhook)

		discordChannels[name] = webhook
	}

	newContent := make(map[string]string, len(discordChannels))
	maps.Copy(newContent, discordChannels)

	path, err := utils.GetPath("json", "discord_channel.json")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to determine file path"})
		return
	}

	if err := utils.WriteJSON(path, newContent); err != nil {
		slog.Error("Failed to write discord_channel.json", "path", path, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save channel configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "channels added successfully"})
}
