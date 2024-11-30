package handlers

import (
	"cartophone-server/internal/pocketbase"
	"cartophone-server/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// AssociateHandler handles associating a card with a playlist
func AssociateHandler(cardDetectedChan <-chan string, modeSwitch chan string, baseURL string, w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] AssociateHandler started. Processing request...")

	// Parse playlist ID
	var payload struct {
		PlaylistID string `json:"playlistId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteResponse(w, http.StatusBadRequest, "Invalid request payload")
		fmt.Println("[DEBUG] Invalid request payload. Exiting handler.")
		return
	}

	if payload.PlaylistID == "" {
		utils.WriteResponse(w, http.StatusBadRequest, "Playlist ID is required")
		fmt.Println("[DEBUG] Playlist ID is missing in the request payload. Exiting handler.")
		return
	}

	fmt.Printf("[DEBUG] Associate mode activated. Waiting for a card to associate with playlist ID: %s\n", payload.PlaylistID)

	select {
	case uid := <-cardDetectedChan:
		fmt.Printf("[DEBUG] Detected card UID in associate mode: %s\n", uid)

		// Check if the card exists in PocketBase
		card, err := pocketbase.CheckCard(baseURL, uid)
		if err != nil {
			utils.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error checking card: %v", err))
			fmt.Printf("[DEBUG] Error checking card in PocketBase: %v\n", err)
			modeSwitch <- ReadMode // Switch back to read mode
			return
		}

		if card != nil && card.PlaylistID != "" {
			if card.PlaylistID == payload.PlaylistID {
				utils.WriteResponse(w, http.StatusConflict, "Card is already associated with this playlist")
				fmt.Printf("[DEBUG] Card %s is already associated with playlist %s\n", uid, payload.PlaylistID)
			} else {
				utils.WriteResponse(w, http.StatusConflict, "Card is already associated with another playlist")
				fmt.Printf("[DEBUG] Card %s is already associated with another playlist\n", uid)
			}
			modeSwitch <- ReadMode // Switch back to read mode
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
			fmt.Printf("[DEBUG] Error adding card to PocketBase: %v\n", err)
			modeSwitch <- ReadMode // Switch back to read mode
			return
		}

		utils.WriteResponse(w, http.StatusOK, fmt.Sprintf("Card %s associated with playlist %s successfully!", uid, payload.PlaylistID))
		fmt.Printf("[DEBUG] Card %s associated with playlist %s successfully. HTTP response sent.\n", uid, payload.PlaylistID)
		modeSwitch <- ReadMode // Switch back to read mode

	case <-time.After(10 * time.Second):
		utils.WriteResponse(w, http.StatusRequestTimeout, "No card detected within 10 seconds")
		fmt.Println("[DEBUG] No card detected within the timeout period")
		modeSwitch <- ReadMode // Switch back to read mode
	}
}