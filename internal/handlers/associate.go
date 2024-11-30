package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"cartophone-server/handlers/utils"
	"cartophone-server/internal/pocketbase"
)

// AssociateHandler handles associating a card with a playlist
func AssociateHandler(cardDetectedChan <-chan string, baseURL string, w http.ResponseWriter, r *http.Request) {
	// Parse playlist ID
	var payload struct {
		PlaylistID string `json:"playlistId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if payload.PlaylistID == "" {
		writeResponse(w, http.StatusBadRequest, "Playlist ID is required")
		return
	}

	fmt.Println("Associate mode activated. Waiting for a card...")

	select {
	case uid := <-cardDetectedChan:
		fmt.Printf("Detected card UID in associate mode: %s\n", uid)

		// Check if the card exists in PocketBase
		card, err := pocketbase.CheckCard(baseURL, uid)
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error checking card: %v", err))
			return
		}

		if card != nil && card.PlaylistID != "" {
			if card.PlaylistID == payload.PlaylistID {
				writeResponse(w, http.StatusConflict, "Card is already associated with this playlist")
			} else {
				writeResponse(w, http.StatusConflict, "Card is already associated with another playlist")
			}
			return
		}

		// Add or update the card
		newCard := pocketbase.Card{UID: uid, PlaylistID: payload.PlaylistID}
		if card == nil {
			err = pocketbase.AddCard(baseURL, newCard)
		} else {
			card.PlaylistID = payload.PlaylistID
			err = pocketbase.UpdateCard(baseURL, *card)
		}
		if err != nil {
			writeResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error updating card: %v", err))
			return
		}

		writeResponse(w, http.StatusOK, fmt.Sprintf("Card %s associated with playlist %s successfully!", uid, payload.PlaylistID))
		return

	case <-time.After(10 * time.Second):
		writeResponse(w, http.StatusRequestTimeout, "No card detected within 10 seconds")
	}
}