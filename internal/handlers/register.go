package handlers

import (
	"fmt"
	"time"
	"net/http"
	"cartophone-server/internal/nfc" // Import the nfc package for register logic
)

// RegisterHandler handles the /register route
func RegisterHandler(cardDetectedChan <-chan string, w http.ResponseWriter, r *http.Request) {
	// In register mode, wait for a card for 10 seconds
	// Create a timeout channel for 10 seconds
	registerTimeout := time.After(10 * time.Second)

	select {
	case uid := <-cardDetectedChan:
		// Card detected within 10 seconds
		fmt.Println("Registering card", uid)
		fmt.Fprintf(w, "Registering card %s\n", uid)
	case <-registerTimeout:
		// No card detected within 10 seconds
		fmt.Println("No card detected")
		fmt.Fprintf(w, "No card detected\n")
	}
}