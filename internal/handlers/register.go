package handlers

import (
	"fmt"
	"net/http"
	"time"
)

// RegisterHandler listens for detected card UID and simulates the "Register" action.
func RegisterHandler(cardDetectedChan <-chan string, w http.ResponseWriter, r *http.Request) {
	// Print that register mode is activated
	fmt.Println("Register mode activated. Waiting for a card...")

	// In register mode, wait for a card for 10 seconds
	registerTimeout := time.After(10 * time.Second)

	select {
	case uid := <-cardDetectedChan:
		// Card detected within 10 seconds
		fmt.Printf("Registering card %s\n", uid)
		fmt.Fprintf(w, "Registering card %s\n", uid)
	case <-registerTimeout:
		// No card detected within timeout
		fmt.Println("No card detected")
		fmt.Fprintf(w, "No card detected\n")
	}
}