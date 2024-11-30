package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"cartophone-server/config"
	"cartophone-server/internal/nfc"
	"cartophone-server/internal/handlers"
)

const (
	ReadMode    = "read"
	AssociateMode = "associate"
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
				} else if mode == AssociateMode {
					fmt.Println("Switched to Associate Mode")
				}

			case uid := <-cardDetectedChan:
				// Handle card detection based on the current mode
				modeLock.Lock()
				if currentMode == ReadMode {
					handlers.HandleReadAction(uid, "http://127.0.0.1:8090")
				} else if currentMode == AssociateMode {
					fmt.Printf("Ignoring card %s because we are in Associate Mode\n", uid)
				}
				modeLock.Unlock()
			}
		}
	}()

	// Start polling for NFC cards
	go reader.StartRead(cardDetectedChan)

	// HTTP endpoint for associate mode
	http.HandleFunc("/associate", func(w http.ResponseWriter, r *http.Request) {
		modeLock.Lock()
		if currentMode == AssociateMode {
			fmt.Println("Already in associate mode")
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Already in associate mode")
			modeLock.Unlock()
			return
		}

		currentMode = AssociateMode
		modeLock.Unlock()

		// Trigger associate mode
		modeSwitch <- AssociateMode

		// Handle the association logic
		handlers.AssociateHandler(cardDetectedChan, "http://127.0.0.1:8090", w, r)

		// Revert back to read mode after association is complete or timeout
		time.Sleep(10 * time.Second)
		modeSwitch <- ReadMode
	})

	// Start the HTTP server
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Main thread remains idle
	select {}
}