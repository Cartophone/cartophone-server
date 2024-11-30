package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cartophone-server/internal/pocketbase"
)

func AssociateHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	// Parse JSON payload
	var payload struct {
		CardID     string `json:"cardId"`
		PlaylistID string `json:"playlistId"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Update the card with the playlistId
	if err := pocketbase.UpdateCard(baseURL, payload.CardID, payload.PlaylistID); err != nil {
		http.Error(w, fmt.Sprintf("Failed to associate card: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Card %s associated with playlist %s successfully!", payload.CardID, payload.PlaylistID)
}