package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

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
