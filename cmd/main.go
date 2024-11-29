package main

import (
	"fmt"
	"log"
	"net/http"

	"cartophone-server/config"      // Import config package
	"cartophone-server/internal/nfc" // Import the internal nfc package
	"cartophone-server/internal/handlers" // Import the internal handlers package
)

func main() {
	// Display a nice start message
	fmt.Println("Cartophone server is starting...")
	fmt.Println("Welcome to Cartophone! Ready to scan NFC cards and interact with Owntone and Pocketbase.")
	fmt.Println("Press Ctrl+C to stop the server.")

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

	// Create a mode switch channel to control the activation of read mode
	modeSwitch := make(chan bool)

	// Start polling for NFC cards (use the device directly from reader)
	go reader.StartRead(cardDetectedChan)

	// Handle the /register endpoint to trigger register mode
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		// Activate register mode, suspend the reading action
		fmt.Println("Register mode activated. Waiting for a card...")
		modeSwitch <- false // Suspend the reading handler

		// Call the RegisterHandler
		handlers.RegisterHandler(cardDetectedChan, w, r)

		// Reactivate the reading mode after registration is done
		modeSwitch <- true // Re-enable the reading handler
	})

	// Start the HTTP server to listen for requests
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Start the read action in a separate goroutine
	go handlers.HandleReadAction(cardDetectedChan, modeSwitch)

	// Main loop to handle detected cards
	for {
		select {
		case uid := <-cardDetectedChan:
			// Print the detected card message
			fmt.Printf("Detected card UID: %s\n", uid)
		}
	}
}