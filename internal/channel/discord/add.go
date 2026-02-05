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
func (h *Handler) Add(c *gin.Context) {
	fn := "DiscordHandler/Add"
	var req struct {
		Datas []struct {
			Name    string `json:"name"`
			Webhook string `json:"webhook"`
		} `json:"datas"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err, fn, "failed to parse request payload")
		return
	}

	if len(req.Datas) == 0 {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "need provide at least one channel data with name and webhook")
		return
	}

	var invalidNames []string
	var invalodWebhooks []string
	for _, data := range req.Datas {
		name := strings.TrimSpace(data.Name)
		webhook := strings.TrimSpace(data.Webhook)

		if !regexName.MatchString(name) {
			slog.Error("invalid channel name format",
				slog.String("channelName", name),
			)
			invalidNames = append(invalidNames, name)
		}
		if !regexWebhook.MatchString(webhook) {
			slog.Error("invalid webhook URL format",
				slog.String("webhook", webhook),
			)
			invalodWebhooks = append(invalodWebhooks, webhook)
		}
	}

	if len(invalidNames) > 0 || len(invalodWebhooks) > 0 {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "invalid channel names or webhook URLs")
		return
	}

	channelsMu.Lock()
	defer channelsMu.Unlock()

	if channels == nil {
		channels = make(map[string]string)
	}

	for _, data := range req.Datas {
		name := strings.TrimSpace(data.Name)
		webhook := strings.TrimSpace(data.Webhook)
		channels[name] = webhook
	}

	newContent := make(map[string]string, len(channels))
	maps.Copy(newContent, channels)

	path, err := utils.GetPath("json", configFile)
	if err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err, fn, "failed to get configuration path")
		return
	}

	if err := utils.WriteJSON(path, newContent); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err, fn, "failed to save configuration")
		return
	}

	c.String(http.StatusOK, fn+": ok")
}
