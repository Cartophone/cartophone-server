package handlers

import (
	"net/http"

	"cartophone-server/internal/owntone"
	"cartophone-server/internal/utils"
)

// PlayerStatusHandler retrieves the status of the OwnTone player
func PlayerStatusHandler(ownToneBaseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.LogToConsole("ERROR", "Invalid request method for PlayerStatusHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	status, err := owntone.GetPlayerStatus(ownToneBaseURL)
	if err != nil {
		utils.LogToConsole("ERROR", "Failed to fetch player status", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.LogToConsole("INFO", "Player status fetched successfully", status)
	utils.WriteJSONResponse(w, http.StatusOK, status)
}