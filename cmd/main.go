package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"cartophone-server/config"
	"cartophone-server/internal/nfc"
	"cartophone-server/internal/handlers"
)

const (
	ReadMode    = "read"
	RegisterMode = "register"
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
	reader, err := nfc.NewReader(config.DevicePath)
	if err != nil {
		log.Fatalf("Failed to initialize NFC reader: %v", err)
	}
	defer reader.Close()

	// Channels for communication
	cardDetectedChan := make(chan string)
	modeSwitch := make(chan string)

	// Synchronization for mode state
	var modeLock sync.Mutex
	currentMode := ReadMode

	// Goroutine to handle NFC reading
	go func() {
		for {
			select {
			case mode := <-modeSwitch:
				modeLock.Lock()
				currentMode = mode
				modeLock.Unlock()
			case uid := <-cardDetectedChan:
				modeLock.Lock()
				if currentMode == ReadMode {
					fmt.Printf("Detected card UID: %s\n", uid)
					handlers.HandleReadAction(uid)
				}
				modeLock.Unlock()
			}
		}
	}()

	// Start polling for NFC cards
	go reader.StartRead(cardDetectedChan)

	// HTTP endpoint for register mode
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		modeLock.Lock()
		if currentMode == RegisterMode {
			fmt.Println("Already in register mode")
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Already in register mode")
			modeLock.Unlock()
			return
		}
		currentMode = RegisterMode
		modeLock.Unlock()

		fmt.Println("Register mode activated. Waiting for a card...")
		modeSwitch <- RegisterMode

		handlers.RegisterHandler(cardDetectedChan, w, r)

		// Revert to read mode after registration
		modeSwitch <- ReadMode
	})

	// Start the HTTP server
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Main loop (blocks indefinitely)
	select {}
}