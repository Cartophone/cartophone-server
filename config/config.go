package config

import (
	"encoding/json"
	"fmt"
	"os"

	"cartophone-server/internal/utils"
)

// Config represents the application's configuration
type Config struct {
	DevicePath     string `json:"devicePath"`
	PocketBaseURL  string `json:"pocket_base_url"`
	OwnToneBaseURL string `json:"owntone_base_url"` // Added OwnTone base URL
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		utils.LogMessage("CONFIG", "Failed to open config file", map[string]interface{}{
			"filePath": filePath,
			"error":    err.Error(),
		})
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		utils.LogMessage("CONFIG", "Failed to decode config file", map[string]interface{}{
			"filePath": filePath,
			"error":    err.Error(),
		})
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Log loaded configuration
	utils.LogMessage("CONFIG", "Configuration loaded successfully", map[string]interface{}{
		"devicePath":     config.DevicePath,
		"pocketBaseURL":  config.PocketBaseURL,
		"ownToneBaseURL": config.OwnToneBaseURL,
	})

	return &config, nil
}