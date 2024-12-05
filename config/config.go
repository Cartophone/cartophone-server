package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the application's configuration
type Config struct {
	DevicePath     string `json:"devicePath"`
	PocketBaseURL  string `json:"pocketbase_url"`
	OwnToneBaseURL string `json:"owntone_base_url"` // Added OwnTone base URL
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	return &config, nil
}