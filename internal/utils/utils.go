package utils

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetPath(arg ...string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		slog.Error("Failed to get working directory", "error", err)
		return "", err
	}

	pathAry := append([]string{wd}, arg...)
	path := filepath.Join(pathAry...)
	path = filepath.ToSlash(path)
	if abs, err := filepath.Abs(path); err == nil {
		path = abs
	}

	return path, nil
}

func GetFile(arg ...string) (map[string]string, error) {
	path, err := GetPath(arg...)
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("file does not exist: %s", path)
		}
		return nil, err
	}

	var result map[string]string
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

func WriteJSON(path string, data map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	dataBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, dataBytes, 0644); err != nil {
		return err
	}
	return nil
}

func ResponseError(c *gin.Context, status int, err error, fn, message string) {
	slog.Error(fn+": "+message,
		slog.Any("error", err),
	)
	c.String(status, message)
}

type ChannelPayload struct {
	Datas []struct {
		Name    string `json:"name"`
		Webhook string `json:"webhook"`
	} `json:"datas"`
}

func CheckChannelPayload(req ChannelPayload, regexName, regexWebhook *regexp.Regexp) error {
	if len(req.Datas) == 0 {
		return fmt.Errorf("need provide at least one channel data with name and webhook")
	}

	var invalidNames []string
	var invalidWebhooks []string
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
			invalidWebhooks = append(invalidWebhooks, webhook)
		}
	}

	if len(invalidNames) > 0 || len(invalidWebhooks) > 0 {
		return fmt.Errorf("invalid channel names or webhook URLs")
	}
	return nil
}
