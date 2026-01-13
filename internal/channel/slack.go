package channel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type SlackRequest struct {
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
	Fields []SlackField `json:"fields,omitempty"`
	Footer *SlackFooter `json:"footer,omitempty"`

	// Bot identity (legacy webhook only)
	Channel   string `json:"channel,omitempty"` // "#channel" or "@user"
	Username  string `json:"username,omitempty"`
	IconEmoji string `json:"icon_emoji,omitempty"` // ":rocket:"
	IconURL   string `json:"icon_url,omitempty"`

	// Threading
	ThreadTS string `json:"thread_ts,omitempty"`
}

type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

type SlackFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

type SlackPayload struct {
	Text        string            `json:"text"`
	Channel     string            `json:"channel,omitempty"`
	Username    string            `json:"username,omitempty"`
	IconEmoji   string            `json:"icon_emoji,omitempty"`
	IconURL     string            `json:"icon_url,omitempty"`
	ThreadTS    string            `json:"thread_ts,omitempty"`
	Attachments []SlackAttachment `json:"attachments,omitempty"`
}

type SlackAttachment struct {
	Fallback   string              `json:"fallback,omitempty"`
	Color      string              `json:"color,omitempty"`
	Pretext    string              `json:"pretext,omitempty"`
	Title      string              `json:"title,omitempty"`
	TitleLink  string              `json:"title_link,omitempty"`
	Text       string              `json:"text,omitempty"`
	ImageURL   string              `json:"image_url,omitempty"`
	ThumbURL   string              `json:"thumb_url,omitempty"`
	Footer     string              `json:"footer,omitempty"`
	FooterIcon string              `json:"footer_icon,omitempty"`
	Timestamp  int64               `json:"ts,omitempty"`
	Fields     []SlackPayloadField `json:"fields,omitempty"`
}

type SlackPayloadField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short,omitempty"`
}

func SendToSlack(request SlackRequest) error {
	payload := SlackPayload{
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
		attachment := SlackAttachment{
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
			fields := make([]SlackPayloadField, 0, len(request.Fields))
			for _, field := range request.Fields {
				fields = append(fields, SlackPayloadField{
					Title: field.Title,
					Value: field.Value,
					Short: field.Short,
				})
			}
			attachment.Fields = fields
		}

		if request.Footer != nil {
			attachment.Footer = request.Footer.Text
			attachment.FooterIcon = request.Footer.IconURL
		}

		payload.Attachments = []SlackAttachment{attachment}
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

func hasAttachment(r SlackRequest) bool {
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
