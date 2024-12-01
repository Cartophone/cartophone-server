package config

import (
	"fmt"
	"os"
	"encoding/json"
)

// Config structure holds the application configuration.
type Config struct {
	DevicePath string `json:"device_path"`
	PocketBaseURL string `json:"pocket_base_url"`
}

// LoadConfig loads configuration from a JSON file.
func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}