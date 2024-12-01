
package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cartophone-server/internal/constants"
	"cartophone-server/internal/pocketbase"
)

func AssociateHandler(cardDetectedChan <-chan string, modeSwitch chan string, baseURL string, w http.ResponseWriter, r *http.Request) {
	fmt.Println("[DEBUG] AssociateHandler started. Processing request...")

	// Parse playlist ID
	var payload struct {
		PlaylistID string `json:"playlistId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		WriteResponse(w, http.StatusBadRequest, "Invalid request payload")
		fmt.Println("[DEBUG] Invalid request payload. Exiting handler.")
		return
	}

	if payload.PlaylistID == "" {
		WriteResponse(w, http.StatusBadRequest, "Playlist ID is required")
		fmt.Println("[DEBUG] Playlist ID is missing in the request payload. Exiting handler.")
		return
	}

	fmt.Printf("[DEBUG] Associate mode requested. Playlist ID: %s\n", payload.PlaylistID)

	// Send AssociateMode signal
	select {
	case modeSwitch <- constants.AssociateMode:
		fmt.Println("[DEBUG] Sent AssociateMode signal to modeSwitch")
	default:
		fmt.Println("[DEBUG] AssociateMode signal already sent")
	}

	// Listen for a card or timeout
	select {
	case uid := <-cardDetectedChan:
		fmt.Printf("[DEBUG] Detected card UID in associate mode: %s\n", uid)

		// Check if the card exists in PocketBase
		card, err := pocketbase.CheckCard(baseURL, uid)
		if err != nil {
			WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error checking card: %v", err))
			fmt.Printf("[DEBUG] Error checking card in PocketBase: %v\n", err)
			return
		}

		if card != nil && card.PlaylistID != "" {
			if card.PlaylistID == payload.PlaylistID {
				WriteResponse(w, http.StatusConflict, "Card is already associated with this playlist")
				fmt.Printf("[DEBUG] Card %s is already associated with playlist %s\n", uid, payload.PlaylistID)
			} else {
				WriteResponse(w, http.StatusConflict, "Card is already associated with another playlist")
				fmt.Printf("[DEBUG] Card %s is already associated with another playlist\n", uid)
			}
			// Switch back to ReadMode after handling
			switchToReadMode(modeSwitch)
			return
		}

		// Add or update the card
		newCard := pocketbase.Card{
			UID:        uid,
			PlaylistID: payload.PlaylistID,
		}
		err = pocketbase.AddCard(baseURL, newCard)
		if err != nil {
			WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error adding card: %v", err))
			fmt.Printf("[DEBUG] Error adding card to PocketBase: %v\n", err)
			// Switch back to ReadMode after handling
			switchToReadMode(modeSwitch)
			return
		}

		WriteResponse(w, http.StatusOK, fmt.Sprintf("Card %s associated with playlist %s successfully!", uid, payload.PlaylistID))
		fmt.Printf("[DEBUG] Card %s associated with playlist %s successfully. HTTP response sent.\n", uid, payload.PlaylistID)

	case <-time.After(10 * time.Second):
		WriteResponse(w, http.StatusRequestTimeout, "No card detected within 10 seconds")
		fmt.Println("[DEBUG] No card detected within the timeout period")
	}

	// Ensure we switch back to ReadMode
	switchToReadMode(modeSwitch)
}

// Helper to switch back to ReadMode
func switchToReadMode(modeSwitch chan string) {
	select {
	case modeSwitch <- constants.ReadMode:
		fmt.Println("[DEBUG] Switching back to Read Mode")
	default:
		fmt.Println("[DEBUG] ReadMode signal already sent")
	}
}