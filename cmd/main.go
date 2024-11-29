package main

import (
	"fmt"
	"log"

	"cartophone-server/internal/nfc" // Import the internal nfc package
)

func main() {
	// Initialize the NFC reader
	reader, err := nfc.NewReader("pn532_i2c:/dev/i2c-1:0x24") // Adjust device path as needed
	if err != nil {
		log.Fatalf("Failed to initialize NFC reader: %v", err)
	}
	defer reader.Close()

	// Create a channel to receive NFC card UIDs asynchronously
	cardDetectedChan := make(chan string)

	// Start polling NFC tags asynchronously
	reader.StartPolling(cardDetectedChan)

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