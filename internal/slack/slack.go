package slack

import (
	"go-notify-hub/internal/utils"
	"log/slog"
	"os"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	slackChannels     map[string]string
	slackChannelsMu   sync.RWMutex
	validSlackWebhook = regexp.MustCompile(`^https://hooks\.slack\.com/services/[A-Z0-9]{8,}/[A-Z0-9]{8,}/[a-zA-Z0-9]{24,}$`)
	validChannelName  = regexp.MustCompile(`^[0-9A-Za-z@_-]+$`)
)

type SlackHandler struct{}

func New() (*SlackHandler, error) {
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

// * GET: /slack/list
func (h *SlackHandler) List(c *gin.Context) {
	slackChannelsMu.RLock()
	defer slackChannelsMu.RUnlock()

	if slackChannels == nil {
		c.JSON(200, gin.H{"channels": map[string]string{}})
		return
	}
	c.JSON(200, gin.H{"channels": slackChannels})
}
