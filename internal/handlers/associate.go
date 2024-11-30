package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cartophone-server/internal/pocketbase"
)

// AssociateHandler handles card association with playlists
func AssociateHandler(cardDetectedChan <-chan string, baseURL string, w http.ResponseWriter, r *http.Request) {
	// Parse the playlist ID from the request
	var payload struct {
		PlaylistID string `json:"playlistId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if payload.PlaylistID == "" {
		http.Error(w, "Playlist ID is required", http.StatusBadRequest)
		return
	}

	fmt.Println("Associate mode activated. Waiting for a card...")

	select {
	case uid := <-cardDetectedChan:
		// Check if the card already exists in PocketBase
		card, err := pocketbase.CheckCard(baseURL, uid)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error checking card: %v", err), http.StatusInternalServerError)
			return
		}

		if card != nil {
			if card.PlaylistID == payload.PlaylistID {
				http.Error(w, "Card is already associated with this playlist", http.StatusConflict)
			} else {
				http.Error(w, "Card is already associated with another playlist", http.StatusConflict)
			}
			return
		}

		// Add the card and associate it with the playlist
		newCard := pocketbase.Card{
			UID:        uid,
			PlaylistID: payload.PlaylistID,
		}

		err = pocketbase.AddCard(baseURL, newCard)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error adding card: %v", err), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Card %s associated with playlist %s successfully!\n", uid, payload.PlaylistID)

	case <-time.After(10 * time.Second):
		// Timeout waiting for card
		http.Error(w, "No card detected within 10 seconds", http.StatusRequestTimeout)
	}
}