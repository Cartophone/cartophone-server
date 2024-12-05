package handlers

import (
	"cartophone-server/internal/pocketbase"
	"cartophone-server/internal/utils"
)

// HandleReadAction handles playing the associated playlist when a card is scanned.
func HandleReadAction(uid string, baseURL string) {
	utils.LogMessage("INFO", "Detected card scanned", map[string]interface{}{"uid": uid})

	// Check if the card exists in PocketBase
	card, err := pocketbase.CheckCard(baseURL, uid)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to check card in PocketBase", err.Error())
		return
	}

	if card == nil {
		utils.LogMessage("INFO", "Card not found in PocketBase", map[string]interface{}{"uid": uid})
		return
	}

	// Fetch the associated playlist
	playlist, err := pocketbase.GetPlaylist(baseURL, card.PlaylistID)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to fetch playlist for card", map[string]interface{}{
			"uid":        uid,
			"playlistId": card.PlaylistID,
			"error":      err.Error(),
		})
		return
	}

	// Simulate playing the playlist
	utils.LogMessage("ACTION", "Playing playlist", map[string]interface{}{
		"playlistName": playlist.Name,
		"uri":          playlist.URI,
		"uid":          uid,
	})
}