package handler

import (
	"goNotify/internal/utils"
	"log/slog"
	"os"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	discordChannels     map[string]string
	discordChannelsMu   sync.RWMutex
	vaildDiscordWebhook = regexp.MustCompile(`^https://discord\.com/api/webhooks/\d{17,20}/[A-Za-z0-9_\-]{68}$`)
	validChannelName    = regexp.MustCompile(`^[0-9A-Za-z@_-]+$`)
)

type DiscordHandler struct{}

func NewDiscordHandler() (*DiscordHandler, error) {
	discordChannelsMu.Lock()
	defer discordChannelsMu.Unlock()

	if discordChannels == nil {
		data, err := utils.GetFile("json", "discord_channel.json")
		if err != nil {
			if os.IsNotExist(err) {
				discordChannels = make(map[string]string)
			} else {
				slog.Error("Failed to read discord_channel.json", "error", err)
				return nil, err
			}
		} else {
			discordChannels = data
		}
	}

	return &DiscordHandler{}, nil
}

// GET: /discord/list
func (h *DiscordHandler) List(c *gin.Context) {
	discordChannelsMu.RLock()
	defer discordChannelsMu.RUnlock()

	if discordChannels == nil {
		c.JSON(200, gin.H{"channels": map[string]string{}})
		return
	}
	c.JSON(200, gin.H{"channels": discordChannels})
}
