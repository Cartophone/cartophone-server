package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"cartophone-server/internal/pocketbase"
	"cartophone-server/internal/utils"
)

// CreateAlarmHandler handles the creation of a new alarm
func CreateAlarmHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	var payload struct {
		PlaylistID string `json:"playlistId"`
		Hour       string `json:"hour"`
		Activated  bool   `json:"activated"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	alarm, err := pocketbase.CreateAlarm(baseURL, payload.PlaylistID, payload.Hour, payload.Activated)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create alarm: %v", err)})
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, alarm)
}

// DeleteAlarmHandler handles the deletion of an alarm by ID
func DeleteAlarmHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	var payload struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err := pocketbase.DeleteAlarm(baseURL, payload.ID)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to delete alarm: %v", err)})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Alarm deleted successfully"})
}

// ListAlarmsHandler handles listing all alarms
func ListAlarmsHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	alarms, err := pocketbase.ListAlarms(baseURL)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to list alarms: %v", err)})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, alarms)
}

// SetAlarmStatusHandler handles updating the activation status of an alarm
func SetAlarmStatusHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	var payload struct {
		ID        string `json:"id"`
		Activated bool   `json:"activated"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err := pocketbase.SetAlarmStatus(baseURL, payload.ID, payload.Activated)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to update alarm status: %v", err)})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Alarm status updated successfully"})
}

// ChangeAlarmPlaylistHandler handles changing the playlist of an alarm
func ChangeAlarmPlaylistHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	var payload struct {
		ID         string `json:"id"`
		PlaylistID string `json:"playlistId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err := pocketbase.ChangeAlarmPlaylist(baseURL, payload.ID, payload.PlaylistID)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to change alarm playlist: %v", err)})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Alarm playlist updated successfully"})
}

// ChangeAlarmHourHandler handles changing the hour of an alarm
func ChangeAlarmHourHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		return
	}

	var payload struct {
		ID   string `json:"id"`
		Hour string `json:"hour"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		return
	}

	err := pocketbase.ChangeAlarmHour(baseURL, payload.ID, payload.Hour)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to change alarm hour: %v", err)})
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Alarm hour updated successfully"})
}