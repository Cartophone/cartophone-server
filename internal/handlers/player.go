package handlers

import (
	"encoding/json"
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

// ListQueueHandler lists the current Owntone player queue
func ListQueueHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.LogMessage("ERROR", "Invalid request method for ListQueueHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	queue, err := owntone.GetQueue(baseURL)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to fetch Owntone queue", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to fetch queue"})
		return
	}

	utils.LogMessage("INFO", "Owntone queue fetched successfully", map[string]interface{}{"count": len(queue)})
	utils.WriteJSONResponse(w, http.StatusOK, queue)
}

// ClearQueueHandler clears the Owntone queue
func ClearQueueHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.LogMessage("ERROR", "Invalid request method for ClearQueueHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	err := owntone.ClearQueue(baseURL)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to clear Owntone queue", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to clear queue"})
		return
	}

	utils.LogMessage("INFO", "Owntone queue cleared successfully", nil)
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Queue cleared successfully"})
}

// AddToQueueHandler adds items to the Owntone queue
func AddToQueueHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.LogMessage("ERROR", "Invalid request method for AddToQueueHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	var payload struct {
		Uris []string `json:"uris"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.LogMessage("ERROR", "Invalid request payload for AddToQueueHandler", nil)
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	if len(payload.Uris) == 0 {
		utils.LogMessage("ERROR", "Empty URIs array in AddToQueueHandler request payload", nil)
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "URIs array cannot be empty"})
		return
	}

	err := owntone.AddToQueue(baseURL, payload.Uris)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to add items to Owntone queue", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to add items to queue"})
		return
	}

	utils.LogMessage("INFO", "Owntone queue items added successfully", map[string]interface{}{"uris": payload.Uris})
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Items added to queue successfully"})
}