package handler

import (
	"goNotify/internal/utils"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

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

	if err := utils.WriteJSON(path, toWrite); err != nil {
		slog.Error("Failed to write discord_channel.json", "path", path, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save channel configuration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "channel deleted successfully"})
}
