package slack

import (
	"maps"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"go-notify-hub/internal/utils"
)

// * POST: /slack/add
// * BODY: { datas: [{ "name": "name", "webhook": "url"}] }
func (h *Handler) Add(c *gin.Context) {
	fn := "SlackHandler/Add"
	var req utils.ChannelPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err, fn, "failed to parse request payload")
		return
	}

	if err := utils.CheckChannelPayload(req, regexName, regexWebhook); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "need provide at least one channel data with name and webhook")
		return
	}

	channelsMu.Lock()
	defer channelsMu.Unlock()

	if channels == nil {
		channels = make(map[string]string)
		return
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
