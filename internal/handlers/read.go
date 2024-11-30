package handlers

import (
	"fmt"
	"cartophone-server/internal/pocketbase"
)

// HandleReadAction checks the card and simulates a playlist action if it exists.
func HandleReadAction(uid string, baseURL string) {
	fmt.Printf("Detected card UID: %s\n", uid)

	card, err := pocketbase.CheckCard(baseURL, uid)
	if err != nil {
		fmt.Printf("Error checking card in PocketBase: %v\n", err)
		return
	}

	if card == nil {
		fmt.Println("Card not found in PocketBase.")
		return
	}

	// Simulate playing the associated playlist
	fmt.Printf("Playing playlist: %s for card UID: %s\n", card.Playlist, uid)
}