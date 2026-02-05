package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-notify-hub/internal/utils"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// * POST: /discord/send/:channelName
func (h *Handler) Send(c *gin.Context) {
	fn := "DiscordHandler/Send"
	channelName := c.Param("channelName")
	if channelName == "" {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "channel name is required")
		return
	}

	if !regexName.MatchString(channelName) {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "invalid channel name format")
		return
	}

	channelsMu.RLock()
	cacheChannels := channels
	channelsMu.RUnlock()

	if cacheChannels == nil {
		c.String(http.StatusOK, fn+": need to add channels first")
		return
	}

	webhook, ok := cacheChannels[channelName]
	if !ok || webhook == "" {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "this channel does not exist")
		return
	}

	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ResponseError(c, http.StatusBadRequest, err, fn, "failed to parse request payload")
		return
	}
	req.WebhookURL = webhook

	if req.Title == "" {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "title is required")
		return
	}

	if req.Description == "" {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "description is required")
		return
	}

	if err := SendMessage(req); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err, fn, "failed to send notification")
		return
	}

	c.String(http.StatusOK, fn+": ok")
}

type Request struct {
	WebhookURL  string `json:"webhook_url"`
	Title       string `json:"title"`
	Description string `json:"description"`

	// Optional
	URL       string       `json:"url,omitempty"`
	Color     string       `json:"color,omitempty"`     // Hex
	Timestamp string       `json:"timestamp,omitempty"` // ISO8601
	Image     string       `json:"image,omitempty"`     // Large
	Thumbnail string       `json:"thumbnail,omitempty"` // Small
	Fields    []EmbedField `json:"fields,omitempty"`
	Footer    *EmbedFooter `json:"footer,omitempty"`
	Author    *EmbedAuthor `json:"author,omitempty"`
	Username  string       `json:"username,omitempty"`
	AvatarURL string       `json:"avatar_url,omitempty"`
}

type Payload struct {
	Username *string `json:"username,omitempty"`
	Avatar   *string `json:"avatar_url,omitempty"`
	Embeds   []Embed `json:"embeds"`
}

type Embed struct {
	Title       string       `json:"title"`
	Description string       `json:"description"`
	URL         string       `json:"url"`
	Color       int          `json:"color"`
	Timestamp   string       `json:"timestamp"`
	Image       *EmbedImage  `json:"image,omitempty"`
	Thumbnail   *EmbedImage  `json:"thumbnail,omitempty"`
	Fields      []EmbedField `json:"fields,omitempty"`
	Footer      *EmbedFooter `json:"footer,omitempty"`
	Author      *EmbedAuthor `json:"author,omitempty"`
}

type EmbedImage struct {
	URL string `json:"url"`
}

type EmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type EmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

type EmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

func SendMessage(request Request) error {
	payload := Payload{
		Embeds: []Embed{
			Embed{
				Title:       request.Title,
				Description: request.Description,
				URL:         request.URL,
				Color:       parseColor(request.Color),
				Timestamp:   parseTimestamp(request.Timestamp),
			},
		},
	}

	if request.Username != "" {
		payload.Username = &request.Username
	}

	if request.AvatarURL != "" {
		payload.Avatar = &request.AvatarURL
	}

	if request.Image != "" {
		payload.Embeds[0].Image = &EmbedImage{
			URL: request.Image,
		}
	}

	if request.Thumbnail != "" {
		payload.Embeds[0].Thumbnail = &EmbedImage{
			URL: request.Thumbnail,
		}
	}

	if len(request.Fields) > 0 {
		fields := []EmbedField{}
		for _, field := range request.Fields {
			fields = append(fields, EmbedField{
				Name:   field.Name,
				Value:  field.Value,
				Inline: field.Inline,
			})
		}
		payload.Embeds[0].Fields = fields
	}

	if request.Footer != nil {
		payload.Embeds[0].Footer = &EmbedFooter{
			Text:    request.Footer.Text,
			IconURL: request.Footer.IconURL,
		}
	}

	if request.Author != nil {
		payload.Embeds[0].Author = &EmbedAuthor{
			Name:    request.Author.Name,
			URL:     request.Author.URL,
			IconURL: request.Author.IconURL,
		}
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	resp, err := http.Post(
		request.WebhookURL,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("discord API error (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

func parseColor(colorStr string) int {
	if colorStr == "" {
		return 0
	}
	var color int
	_, err := fmt.Sscanf(colorStr, "#%06x", &color)
	if err != nil {
		return 0
	}
	return color
}

func parseTimestamp(timestamp string) string {
	if timestamp == "" {
		return time.Now().Format(time.RFC3339)
	}
	return timestamp
}
