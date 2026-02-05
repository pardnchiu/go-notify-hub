package discord

import (
	"go-notify-hub/internal/utils"
	"maps"
	"net/http"

	"github.com/gin-gonic/gin"
)

// * DELETE: /discord/delete/:channelName
func (h *Handler) Delete(c *gin.Context) {
	fn := "DiscordHandler/Delete"
	channelName := c.Param("channelName")
	if channelName == "" {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "channel name is required")
		return
	}

	if !regexName.MatchString(channelName) {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "invalid channel name format")
		return
	}

	channelsMu.Lock()
	defer channelsMu.Unlock()

	if channels == nil {
		c.String(http.StatusOK, fn+": need to add channels first")
		return
	}

	delete(channels, channelName)
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
