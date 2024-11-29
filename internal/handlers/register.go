package handlers

import (
	"fmt"
	"net/http"
	"time"
)

// RegisterHandler handles the /register route
func RegisterHandler(cardDetectedChan <-chan string, w http.ResponseWriter, r *http.Request) {
	// Waiting for a card to be detected within 10 seconds
	select {
	case uid := <-cardDetectedChan:
		// Card detected, simulate registration action
		fmt.Printf("Registering card %s\n", uid)
		fmt.Fprintf(w, "Registering card %s\n", uid)
	case <-time.After(10 * time.Second):
		// No card detected within timeout
		fmt.Println("No card detected")
		fmt.Fprintf(w, "No card detected\n")
	}
}