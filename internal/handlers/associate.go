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
    // Parse playlist ID
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

    // Wait for card or timeout
    select {
    case uid := <-cardDetectedChan:
        fmt.Printf("Detected card UID in associate mode: %s\n", uid)
        card, err := pocketbase.CheckCard(baseURL, uid)
        if err != nil {
            http.Error(w, fmt.Sprintf("Error checking card: %v", err), http.StatusInternalServerError)
            return
        }

        if card != nil && card.PlaylistID != "" {
            if card.PlaylistID == payload.PlaylistID {
                http.Error(w, "Card is already associated with this playlist", http.StatusConflict)
            } else {
                http.Error(w, "Card is already associated with another playlist", http.StatusConflict)
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
            http.Error(w, fmt.Sprintf("Error updating card: %v", err), http.StatusInternalServerError)
            return
        }

        fmt.Fprintf(w, "Card %s associated with playlist %s successfully!\n", uid, payload.PlaylistID)

    case <-time.After(10 * time.Second):
        fmt.Println("No card detected within the timeout period")
        http.Error(w, "No card detected within 10 seconds", http.StatusRequestTimeout)
    }
}