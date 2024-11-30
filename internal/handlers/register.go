package handlers

import (
	"fmt"
	"net/http"
	"time"

	"cartophone-server/internal/pocketbase"
)

// RegisterHandler adds a new card to the PocketBase database
func RegisterHandler(cardDetectedChan <-chan string, baseURL string, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Register mode activated. Waiting for a card...")

	// Wait for a card for 10 seconds
	registerTimeout := time.After(10 * time.Second)

	select {
	case uid := <-cardDetectedChan:
		// Card detected within timeout
		fmt.Printf("Registering card %s\n", uid)

		card := pocketbase.Card{
			UID:        uid,
			PlaylistID: "", // Leave PlaylistID empty during registration
		}

		err := pocketbase.AddCard(baseURL, card)
		if err != nil {
			fmt.Printf("Error registering card: %v\n", err)
			http.Error(w, "Failed to register card", http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Card %s registered successfully!\n", uid)
	case <-registerTimeout:
		// No card detected within timeout
		fmt.Println("No card detected")
		fmt.Fprintf(w, "No card detected\n")
	}
}