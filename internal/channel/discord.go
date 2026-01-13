package channel

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DiscordRequest struct {
	WebhookURL  string `json:"webhook_url"`
	Title       string `json:"title"`
	Description string `json:"description"`

	// Optional
	URL       string              `json:"url,omitempty"`
	Color     string              `json:"color,omitempty"`     // Hex
	Timestamp string              `json:"timestamp,omitempty"` // ISO8601
	Image     string              `json:"image,omitempty"`     // Large
	Thumbnail string              `json:"thumbnail,omitempty"` // Small
	Fields    []DiscordEmbedField `json:"fields,omitempty"`
	Footer    *DiscordEmbedFooter `json:"footer,omitempty"`
	Author    *DiscordEmbedAuthor `json:"author,omitempty"`
	Username  string              `json:"username,omitempty"`
	AvatarURL string              `json:"avatar_url,omitempty"`
}

type DiscordPayload struct {
	Username *string        `json:"username,omitempty"`
	Avatar   *string        `json:"avatar_url,omitempty"`
	Embeds   []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
	Title       string              `json:"title"`
	Description string              `json:"description"`
	URL         string              `json:"url"`
	Color       int                 `json:"color"`
	Timestamp   string              `json:"timestamp"`
	Image       *DiscordEmbedImage  `json:"image,omitempty"`
	Thumbnail   *DiscordEmbedImage  `json:"thumbnail,omitempty"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Footer      *DiscordEmbedFooter `json:"footer,omitempty"`
	Author      *DiscordEmbedAuthor `json:"author,omitempty"`
}

type DiscordEmbedImage struct {
	URL string `json:"url"`
}

type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

type DiscordEmbedFooter struct {
	Text    string `json:"text"`
	IconURL string `json:"icon_url,omitempty"`
}

type DiscordEmbedAuthor struct {
	Name    string `json:"name"`
	URL     string `json:"url,omitempty"`
	IconURL string `json:"icon_url,omitempty"`
}

func SendToDiscord(request DiscordRequest) error {
	payload := DiscordPayload{
		Embeds: []DiscordEmbed{
			DiscordEmbed{
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
		payload.Embeds[0].Image = &DiscordEmbedImage{
			URL: request.Image,
		}
	}

	if request.Thumbnail != "" {
		payload.Embeds[0].Thumbnail = &DiscordEmbedImage{
			URL: request.Thumbnail,
		}
	}

	if len(request.Fields) > 0 {
		fields := []DiscordEmbedField{}
		for _, field := range request.Fields {
			fields = append(fields, DiscordEmbedField{
				Name:   field.Name,
				Value:  field.Value,
				Inline: field.Inline,
			})
		}
		payload.Embeds[0].Fields = fields
	}

	if request.Footer != nil {
		payload.Embeds[0].Footer = &DiscordEmbedFooter{
			Text:    request.Footer.Text,
			IconURL: request.Footer.IconURL,
		}
	}

	if request.Author != nil {
		payload.Embeds[0].Author = &DiscordEmbedAuthor{
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
