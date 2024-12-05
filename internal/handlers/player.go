package handlers

import (
	"net/http"

	"cartophone-server/internal/owntone"
	"cartophone-server/internal/utils"
)

// PlayerStatusHandler retrieves the status of the OwnTone player
func PlayerStatusHandler(ownToneBaseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.LogMessage("ERROR", "Invalid request method for PlayerStatusHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	status, err := owntone.GetPlayerStatus(ownToneBaseURL)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to fetch player status", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.LogMessage("INFO", "Player status fetched successfully", status)
	utils.WriteJSONResponse(w, http.StatusOK, status)
}

// PlayHandler triggers the play action on the Owntone player
func PlayHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.LogMessage("ERROR", "Invalid request method for PlayHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	utils.LogMessage("INFO", "Received request to play", nil)

	err := owntone.Play(baseURL)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to play Owntone", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.LogMessage("INFO", "Playback started successfully", nil)
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Playback started"})
}

// PauseHandler triggers the pause action on the Owntone player
func PauseHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.LogMessage("ERROR", "Invalid request method for PauseHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	utils.LogMessage("INFO", "Received request to pause", nil)

	err := owntone.Pause(baseURL)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to pause Owntone", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.LogMessage("INFO", "Playback paused successfully", nil)
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Playback paused"})
}