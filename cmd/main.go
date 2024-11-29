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

	// Start polling for NFC tags
	reader.StartPolling() // Call the function that continuously scans for NFC tags
}