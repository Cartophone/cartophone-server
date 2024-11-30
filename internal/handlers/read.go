package handlers

import (
	"fmt"

	"cartophone-server/internal/pocketbase"
)

// HandleReadAction plays the associated playlist when a card is scanned
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

	// Fetch the associated playlist
	playlist, err := pocketbase.GetPlaylist(baseURL, card.PlaylistId)
	if err != nil {
		fmt.Printf("Error fetching playlist: %v\n", err)
		return
	}

	// Simulate playing the playlist
	fmt.Printf("Playing playlist '%s' with URI: %s\n", playlist.Name, playlist.URI)
}