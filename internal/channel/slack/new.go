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
	configFile   = "slack_channel.json"
	channels     map[string]string
	channelsMu   sync.RWMutex
	regexName    = regexp.MustCompile(`^[0-9A-Za-z@_-]+$`)
	regexWebhook = regexp.MustCompile(`^https://hooks\.slack\.com/services/[A-Z0-9]{8,}/[A-Z0-9]{8,}/[a-zA-Z0-9]{24,}$`)
)

type Handler struct{}

func New() (*Handler, error) {
	channelsMu.Lock()
	defer channelsMu.Unlock()

	if channels == nil {
		data, err := utils.GetFile("json", configFile)
		if err != nil {
			if os.IsNotExist(err) {
				channels = make(map[string]string)
			} else {
				slog.Error("failed to get configuration path",
					slog.Any("error", err))
				return nil, err
			}
		} else {
			channels = data
		}
	}

	return &Handler{}, nil
}

// * GET: /slack/list
func (h *Handler) List(c *gin.Context) {
	channelsMu.RLock()
	defer channelsMu.RUnlock()

	if channels == nil {
		c.JSON(200, gin.H{"channels": map[string]string{}})
		return
	}
	c.JSON(200, gin.H{"channels": channels})
}
