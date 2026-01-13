package handler

import (
	"goNotify/internal/utils"
	"log/slog"
	"os"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	slackChannels   map[string]string
	slackChannelsMu sync.RWMutex
)

type SlackHandler struct{}

func NewSlackHandler() (*SlackHandler, error) {
	slackChannelsMu.Lock()
	defer slackChannelsMu.Unlock()

	if slackChannels == nil {
		data, err := utils.GetFile("json", "slack_channel.json")
		if err != nil {
			if os.IsNotExist(err) {
				slackChannels = make(map[string]string)
			} else {
				slog.Error("Failed to read slack_channel.json", "error", err)
				return nil, err
			}
		} else {
			slackChannels = data
		}
	}

	return &SlackHandler{}, nil
}

// GET: /slack/list
func (h *SlackHandler) List(c *gin.Context) {
	slackChannelsMu.RLock()
	defer slackChannelsMu.RUnlock()

	if slackChannels == nil {
		c.JSON(200, gin.H{"channels": map[string]string{}})
		return
	}
	c.JSON(200, gin.H{"channels": slackChannels})
}
