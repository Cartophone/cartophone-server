package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"cartophone-server/internal/pocketbase"
	"cartophone-server/internal/utils"
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

	select {
	case uid := <-cardDetectedChan:
		fmt.Printf("Detected card UID in associate mode: %s\n", uid)

		card, err := pocketbase.CheckCard(baseURL, uid)
		if err != nil {
			utils.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error checking card: %v", err))
			return
		}

		if card != nil {
			if card.PlaylistID == "" {
				card.PlaylistID = payload.PlaylistID
				err = pocketbase.UpdateCard(baseURL, *card)
				if err != nil {
					utils.WriteResponse(w, http.StatusInternalServerError, fmt.Sprintf("Error updating card: %v", err))
					return
				}
				utils.WriteResponse(w, http.StatusOK, fmt.Sprintf("Card %s associated with playlist %s successfully!", uid, payload.PlaylistID))
				return
			} else if card.PlaylistID == payload.PlaylistID {
				utils.WriteResponse(w, http.StatusConflict, "Card is already associated with this playlist")
				return
			} else {
				utils.WriteResponse(w, http.StatusConflict, "Card is already associated with another playlist")
				return
			}
		}

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

	case <-time.After(10 * time.Second):
		utils.WriteResponse(w, http.StatusRequestTimeout, "No card detected within 10 seconds")
	}
}