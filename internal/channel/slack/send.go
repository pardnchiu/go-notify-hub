package slack

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

// * POST: /slack/send/:channelName
func (h *Handler) Send(c *gin.Context) {
	fn := "SlackHandler/Send"
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

	if req.Text == "" {
		utils.ResponseError(c, http.StatusBadRequest, nil, fn, "text is required")
		return
	}

	if err := send(req); err != nil {
		utils.ResponseError(c, http.StatusInternalServerError, err, fn, "failed to send notification")
		return
	}

	c.String(http.StatusOK, fn+": ok")
}

type Request struct {
	WebhookURL string `json:"webhook_url"`
	Text       string `json:"text"` // Required: fallback & notification text

	// Optional content (用 Attachment 呈現)
	Title       string `json:"title,omitempty"`
	TitleLink   string `json:"title_link,omitempty"`
	Description string `json:"description,omitempty"` // Attachment text
	Pretext     string `json:"pretext,omitempty"`     // Text above attachment

	// Optional styling
	Color     string `json:"color,omitempty"`     // Hex "#FF5733" or "good"/"warning"/"danger"
	Timestamp int64  `json:"timestamp,omitempty"` // Unix timestamp
	Image     string `json:"image,omitempty"`     // Large image
	Thumbnail string `json:"thumbnail,omitempty"` // Thumb image (右側)

	// Optional extras
	Fields []Field `json:"fields,omitempty"`
	Footer *Footer `json:"footer,omitempty"`

	// Bot identity (legacy webhook only)
	Channel   string `json:"channel,omitempty"` // "#channel" or "@user"
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"` // ":rocket:"
	IconURL   string `json:"icon_url,omitempty"`

	// Threading
	ThreadTS string `json:"thread_ts,omitempty"`
}

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

type Footer struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

type Payload struct {
	Text        string       `json:"text"`
	Channel     string       `json:"channel,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	ThreadTS    string       `json:"thread_ts,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Fallback   string  `json:"fallback,omitempty"`
	Color      string  `json:"color,omitempty"`
	Pretext    string  `json:"pretext,omitempty"`
	Title      string  `json:"title,omitempty"`
	TitleLink  string  `json:"title_link,omitempty"`
	Text       string  `json:"text,omitempty"`
	ImageURL   string  `json:"image_url,omitempty"`
	ThumbURL   string  `json:"thumb_url,omitempty"`
	Footer     string  `json:"footer,omitempty"`
	FooterIcon string  `json:"footer_icon,omitempty"`
	Timestamp  int64   `json:"ts,omitempty"`
	Fields     []Field `json:"fields,omitempty"`
}

func send(request Request) error {
	payload := Payload{
		Text: request.Text,
	}

	if request.Channel != "" {
		payload.Channel = request.Channel
	}
	if request.Username != "" {
		payload.Username = request.Username
	}
	if request.IconEmoji != "" {
		payload.IconEmoji = request.IconEmoji
	}
	if request.IconURL != "" {
		payload.IconURL = request.IconURL
	}
	if request.ThreadTS != "" {
		payload.ThreadTS = request.ThreadTS
	}

	if hasAttachment(request) {
		attachment := Attachment{
			Fallback:  request.Text,
			Color:     parseSlackColor(request.Color),
			Pretext:   request.Pretext,
			Title:     request.Title,
			TitleLink: request.TitleLink,
			Text:      request.Description,
			Timestamp: parseSlackTimestamp(request.Timestamp),
		}

		if request.Image != "" {
			attachment.ImageURL = request.Image
		}

		if request.Thumbnail != "" {
			attachment.ThumbURL = request.Thumbnail
		}

		if len(request.Fields) > 0 {
			attachment.Fields = request.Fields
		}

		if request.Footer != nil {
			attachment.Footer = request.Footer.Text
			attachment.FooterIcon = request.Footer.IconURL
		}

		payload.Attachments = []Attachment{attachment}
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

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("slack API error (status %d): %s", resp.StatusCode, string(body))
	}

	if string(body) != "ok" {
		return fmt.Errorf("slack API error: %s", string(body))
	}

	return nil
}

func hasAttachment(r Request) bool {
	return r.Title != "" ||
		r.Description != "" ||
		r.Pretext != "" ||
		r.Color != "" ||
		r.Image != "" ||
		r.Thumbnail != "" ||
		len(r.Fields) > 0 ||
		r.Footer != nil
}

func parseSlackColor(colorStr string) string {
	if colorStr == "" {
		return ""
	}

	switch colorStr {
	case "good", "warning", "danger":
		return colorStr
	default:
		return colorStr
	}
}

func parseSlackTimestamp(timestamp int64) int64 {
	if timestamp == 0 {
		return time.Now().Unix()
	}
	return timestamp
}
