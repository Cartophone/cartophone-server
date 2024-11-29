package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

	// Start polling NFC tags asynchronously, passing the channel
	go reader.StartPolling(cardDetectedChan)

	// Handle the /register endpoint to activate register mode
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		// Register mode: wait for a card for 10 seconds
		fmt.Println("Register mode activated. Waiting for a card...")

		select {
		case uid := <-cardDetectedChan:
			// Card detected within 10 seconds
			fmt.Println(uid)
			fmt.Fprintf(w, "%s\n", uid)
		case <-time.After(10 * time.Second):
			// No card detected within 10 seconds
			fmt.Println("No card detected")
			fmt.Fprintf(w, "No card detected\n")
		}
	})

	// Start the HTTP server to listen for requests
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Main loop to handle detected cards (continuously running)
	for {
		select {
		case uid := <-cardDetectedChan:
			// Print the detected card message
			fmt.Println(uid)
		}
	}
}