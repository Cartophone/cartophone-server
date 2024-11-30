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

	// Create channels for NFC card detection and mode switching
	cardDetectedChan := make(chan string)
	modeSwitch := make(chan string)

	// Synchronization for mode state
	var modeLock sync.Mutex
	currentMode := ReadMode

	// Goroutine to manage NFC card detection and mode state
	go func() {
		for {
			select {
			case mode := <-modeSwitch:
				// Update the mode state
				modeLock.Lock()
				currentMode = mode
				modeLock.Unlock()

				if mode == ReadMode {
					fmt.Println("Switched to Read Mode")
				} else if mode == RegisterMode {
					fmt.Println("Switched to Register Mode")
				}

			case uid := <-cardDetectedChan:
				// Handle card detection based on the current mode
				modeLock.Lock()
				if currentMode == ReadMode {
					handlers.HandleReadAction(uid, "http://127.0.0.1:8090")
				} else if currentMode == RegisterMode {
					fmt.Printf("Ignoring card %s because we are in Register Mode\n", uid)
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

		// Trigger register mode
		modeSwitch <- RegisterMode

		// Handle the registration logic
		handlers.RegisterHandler(cardDetectedChan, "http://127.0.0.1:8090", w, r)

		// Revert back to read mode after registration is complete
		modeSwitch <- ReadMode
	})

	// Start the HTTP server
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Main thread remains idle
	select {}
}