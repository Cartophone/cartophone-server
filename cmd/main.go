package main

import (
	"fmt"
	"log"

	"cartophone-server/config"   // Import config package
	"cartophone-server/internal/nfc" // Import the internal nfc package
)

func main() {
	// Load the configuration from config.json
	config, err := config.LoadConfig("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize the NFC reader using the device path from the config
	reader, err := nfc.NewReader(config.DevicePath) // Use the loaded DevicePath
	if err != nil {
		log.Fatalf("Failed to initialize NFC reader: %v", err)
	}
	defer reader.Close()

	// Start polling for NFC tags
	reader.StartPolling() // Call the function that continuously scans for NFC tags
}