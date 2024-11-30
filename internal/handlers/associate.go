package handlers

import (
	"cartophone-server/internal/pocketbase"
	"cartophone-server/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

func AssociateHandler(cardDetectedChan <-chan string, baseURL string, w http.ResponseWriter, r *http.Request) {
	// Parse playlist ID
	var payload struct {
		PlaylistID string `json:"playlistId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if payload.PlaylistID == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Playlist ID is required")
		return
	}

	fmt.Println("Associate mode activated. Waiting for a card...")

	// Use a WaitGroup to ensure response is handled exactly once
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		select {
		case uid := <-cardDetectedChan:
			fmt.Printf("Detected card UID in associate mode: %s\n", uid)

			// Check if the card exists in PocketBase
			card, err := pocketbase.CheckCard(baseURL, uid)
			if err != nil {
				utils.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error checking card: %v", err))
				return
			}

			if card != nil && card.PlaylistID != "" {
				if card.PlaylistID == payload.PlaylistID {
					utils.WriteResponse(w, http.StatusConflict, "Card is already associated with this playlist")
				} else {
					utils.WriteResponse(w, http.StatusConflict, "Card is already associated with another playlist")
				}
				return
			}

			// Add or update the card
			newCard := pocketbase.Card{
				UID:        uid,
				PlaylistID: payload.PlaylistID,
			}
			err = pocketbase.AddCard(baseURL, newCard)
			if err != nil {
				utils.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error adding card: %v", err))
				return
			}

			utils.WriteResponse(w, http.StatusOK, fmt.Sprintf("Card %s associated with playlist %s successfully!", uid, payload.PlaylistID))
			return

		case <-time.After(10 * time.Second):
			utils.WriteResponse(w, http.StatusRequestTimeout, "No card detected within 10 seconds")
			return
		}
	}()

	// Wait for the goroutine to finish before exiting the handler
	wg.Wait()
}