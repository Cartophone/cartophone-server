package main

import (
	"fmt"
	"log"
	"time"

	"cartophone-server/internal/nfc" // Import internal NFC package
)

func main() {
	// Initialize the NFC reader
	reader, err := nfc.NewReader("pn532_i2c:/dev/i2c-1:0x24") // Adjust based on your device
	if err != nil {
		log.Fatalf("Failed to initialize NFC reader: %v", err)
	}
	defer reader.Close()

	// Start polling for NFC tags
	fmt.Println("NFC reader initialized. Scanning for NFC tags...")

	for {
		// Define modulation types for polling (adjust this as needed)
		modulations := []nfc.Modulation{
			{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
		}

		// Poll for a target (NFC card/tag)
		uid, err := reader.Scan(modulations, 10, 300*time.Millisecond) // 10 attempts, 300ms period
		if err != nil {
			log.Printf("Error scanning NFC tag: %v", err)
			continue
		}

		if uid != "" {
			// Print the detected NFC tag's UID
			fmt.Printf("Tag detected! UID: %s\n", uid)
		} else {
			// No tag detected within the polling period
			fmt.Println("No NFC tag detected.")
		}

		// Wait before polling again
		time.Sleep(1 * time.Second)
	}
}