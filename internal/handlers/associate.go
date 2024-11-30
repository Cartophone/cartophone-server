package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cartophone-server/internal/pocketbase"
)

// AssociateHandler handles associating a card with a playlist
func AssociateHandler(cardDetectedChan <-chan string, baseURL string, w http.ResponseWriter, r *http.Request) {
	// Parse the playlist ID from the request body
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

	// Use a separate channel to notify when a response is sent
	responseSent := make(chan bool, 1)

	go func() {
		select {
		case uid := <-cardDetectedChan:
			// Card detected, handle association logic
			card, err := pocketbase.CheckCard(baseURL, uid)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error checking card: %v", err), http.StatusInternalServerError)
				responseSent <- true
				return
			}

			if card != nil {
				if card.PlaylistID == "" {
					// Update the card with the new PlaylistID
					card.PlaylistID = payload.PlaylistID
					err = pocketbase.UpdateCard(baseURL, *card)
					if err != nil {
						http.Error(w, fmt.Sprintf("Error updating card: %v", err), http.StatusInternalServerError)
						responseSent <- true
						return
					}
					fmt.Fprintf(w, "Card %s associated with playlist %s successfully!\n", uid, payload.PlaylistID)
					responseSent <- true
					return
				} else if card.PlaylistID == payload.PlaylistID {
					http.Error(w, "Card is already associated with this playlist", http.StatusConflict)
					responseSent <- true
					return
				} else {
					http.Error(w, "Card is already associated with another playlist", http.StatusConflict)
					responseSent <- true
					return
				}
			}

			// If the card does not exist, create a new one and associate it with the playlist
			newCard := pocketbase.Card{
				UID:        uid,
				PlaylistID: payload.PlaylistID,
			}
			err = pocketbase.AddCard(baseURL, newCard)
			if err != nil {
				http.Error(w, fmt.Sprintf("Error adding card: %v", err), http.StatusInternalServerError)
				responseSent <- true
				return
			}

			fmt.Fprintf(w, "Card %s associated with playlist %s successfully!\n", uid, payload.PlaylistID)
			responseSent <- true

		default:
			// Check responseSent to avoid blocking
			select {
			case <-responseSent:
				return
			default:
			}
		}
	}()

	// Wait for up to 10 seconds for a response
	select {
	case <-responseSent:
		// Response was handled in the goroutine
		return
	case <-time.After(10 * time.Second):
		// Timeout, no card detected
		http.Error(w, "No card detected within 10 seconds", http.StatusRequestTimeout)
		return
	}
}