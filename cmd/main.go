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

	// Create a channel to receive NFC card UIDs asynchronously
	cardDetectedChan := make(chan string)

	// Start polling NFC tags asynchronously, passing the channel to StartPolling
	go reader.StartPolling(cardDetectedChan)

	// Main loop to handle detected cards
	for {
		select {
		case uid := <-cardDetectedChan:
			// Trigger actions when a card is detected
			fmt.Printf("Card detected! UID: %s\n", uid)
			// You can add further logic here, e.g., interacting with Pocketbase or Owntone.
		}
	}
}