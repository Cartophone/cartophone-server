package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"cartophone-server/internal/constants"
	"cartophone-server/internal/pocketbase"
	"cartophone-server/internal/utils"
)

func AssociateCardHandler(cardDetectedChan <-chan string, modeSwitch chan string, baseURL string, w http.ResponseWriter, r *http.Request) {
	utils.LogMessage("DEBUG", "AssociateHandler started. Processing request...", nil)

	// Parse playlist ID and replaceCard flag
	var payload struct {
		PlaylistID  string `json:"playlistId"`
		ReplaceCard bool   `json:"replaceCard,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		utils.LogMessage("ERROR", "Invalid request payload", err.Error())
		return
	}

	if payload.PlaylistID == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Playlist ID is required"})
		utils.LogMessage("ERROR", "Playlist ID is missing in the request payload", nil)
		return
	}

	utils.LogMessage("DEBUG", "Associate mode requested", map[string]interface{}{
		"playlistId":  payload.PlaylistID,
		"replaceCard": payload.ReplaceCard,
	})

	// Send AssociateMode signal
	select {
	case modeSwitch <- constants.AssociateMode:
		utils.LogMessage("DEBUG", "Sent AssociateMode signal to modeSwitch", nil)
	default:
		utils.LogMessage("DEBUG", "AssociateMode signal already sent", nil)
	}

	// Listen for a card or timeout
	select {
	case uid := <-cardDetectedChan:
		utils.LogMessage("DEBUG", "Detected card UID in associate mode", uid)

		// Check if the card exists in PocketBase
		card, err := pocketbase.CheckCard(baseURL, uid)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error checking card in PocketBase"})
			utils.LogMessage("ERROR", "Error checking card in PocketBase", err.Error())
			return
		}

		if card != nil && card.PlaylistID != "" {
			if card.PlaylistID == payload.PlaylistID {
				utils.WriteJSONResponse(w, http.StatusConflict, map[string]string{
					"message": "Card is already associated with this playlist",
				})
				utils.LogMessage("INFO", "Card is already associated with the requested playlist", card)
			} else if payload.ReplaceCard {
				card.PlaylistID = payload.PlaylistID
				err = pocketbase.UpdateCard(baseURL, *card)
				if err != nil {
					utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error updating card in PocketBase"})
					utils.LogMessage("ERROR", "Error updating card in PocketBase", err.Error())
					switchToReadMode(modeSwitch)
					return
				}

				utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
					"message":    "Card reassigned to the new playlist",
					"cardId":     card.ID,
					"playlistId": card.PlaylistID,
				})
				utils.LogMessage("INFO", "Card reassigned to the new playlist", card)
			} else {
				utils.WriteJSONResponse(w, http.StatusConflict, map[string]interface{}{
					"message":    "Card is already associated with another playlist",
					"cardId":     card.ID,
					"playlistId": card.PlaylistID,
				})
				utils.LogMessage("INFO", "Card is associated with another playlist", card)
			}
			switchToReadMode(modeSwitch)
			return
		}

		// Add a new card
		newCard := pocketbase.Card{
			UID:        uid,
			PlaylistID: payload.PlaylistID,
		}
		err = pocketbase.AddCard(baseURL, newCard)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Error adding card to PocketBase"})
			utils.LogMessage("ERROR", "Error adding card to PocketBase", err.Error())
			switchToReadMode(modeSwitch)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, map[string]string{
			"message":    "Card associated successfully",
			"cardId":     newCard.ID,
			"playlistId": newCard.PlaylistID,
		})
		utils.LogMessage("INFO", "Card associated successfully", newCard)

	case <-time.After(10 * time.Second):
		utils.WriteJSONResponse(w, http.StatusRequestTimeout, map[string]string{"error": "No card detected within 10 seconds"})
		utils.LogMessage("DEBUG", "No card detected within timeout period", nil)
	}

	// Ensure we switch back to ReadMode
	switchToReadMode(modeSwitch)
}

func switchToReadMode(modeSwitch chan string) {
	select {
	case modeSwitch <- constants.ReadMode:
		utils.LogMessage("DEBUG", "Switched back to Read Mode", nil)
	default:
		utils.LogMessage("DEBUG", "ReadMode signal already sent", nil)
	}
}