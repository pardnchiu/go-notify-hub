package discord

import (
	"go-notify-hub/internal/utils"
	"log/slog"
	"maps"
	"net/http"

	"github.com/gin-gonic/gin"
)

// * DELETE: /discord/delete/:channelName
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

	discordChannelsMu.Lock()
	defer discordChannelsMu.Unlock()

	if discordChannels == nil {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
		return
	}

	delete(discordChannels, channelName)
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

	c.JSON(http.StatusOK, gin.H{"message": "channel deleted successfully"})
}
