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

// QueueListHandler retrieves the Owntone player queue
func QueueListHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.LogMessage("ERROR", "Invalid request method for QueueListHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	queue, err := owntone.GetQueue(baseURL)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to fetch queue", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.LogMessage("INFO", "Owntone queue fetched successfully", nil)
	utils.WriteJSONResponse(w, http.StatusOK, queue)
}

// QueueClearHandler clears the Owntone player queue
func QueueClearHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.LogMessage("ERROR", "Invalid request method for QueueClearHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	err := owntone.ClearQueue(baseURL)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to clear queue", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.LogMessage("INFO", "Owntone queue cleared successfully", nil)
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Queue cleared successfully"})
}

// QueueAddHandler adds a track to the Owntone player queue
func QueueAddHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.LogMessage("ERROR", "Invalid request method for QueueAddHandler", map[string]string{"method": r.Method})
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	var payload struct {
		TrackURI string `json:"uri"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.LogMessage("ERROR", "Invalid request payload for QueueAddHandler", nil)
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	if payload.TrackURI == "" {
		utils.LogMessage("ERROR", "Missing URI in QueueAddHandler", nil)
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "URI is required"})
		return
	}

	err := owntone.AddToQueue(baseURL, payload.TrackURI)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to add track to queue", map[string]string{"error": err.Error()})
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	utils.LogMessage("INFO", "Track added to Owntone queue successfully", map[string]string{"uri": payload.TrackURI})
	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Track added to queue successfully"})
}