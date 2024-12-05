package handlers

import (
	"encoding/json"
	"net/http"

	"cartophone-server/internal/pocketbase"
	"cartophone-server/internal/utils"
)

// CreateAlarmHandler handles the creation of a new alarm
func CreateAlarmHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		utils.LogMessage("ERROR", "Invalid request method for CreateAlarmHandler", nil)
		return
	}

	var payload struct {
		PlaylistID string `json:"playlistId"`
		Hour       string `json:"hour"`
		Activated  bool   `json:"activated"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		utils.LogMessage("ERROR", "Failed to decode request body for CreateAlarmHandler", err.Error())
		return
	}

	alarm, err := pocketbase.CreateAlarm(baseURL, payload.PlaylistID, payload.Hour, payload.Activated)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to create alarm"})
		utils.LogMessage("ERROR", "Failed to create alarm in PocketBase", err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusCreated, alarm)
	utils.LogMessage("INFO", "Alarm created successfully", alarm)
}

// DeleteAlarmHandler handles the deletion of an alarm by ID
func DeleteAlarmHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		utils.LogMessage("ERROR", "Invalid request method for DeleteAlarmHandler", nil)
		return
	}

	var payload struct {
		ID string `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		utils.LogMessage("ERROR", "Failed to decode request body for DeleteAlarmHandler", err.Error())
		return
	}

	err := pocketbase.DeleteAlarm(baseURL, payload.ID)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to delete alarm"})
		utils.LogMessage("ERROR", "Failed to delete alarm in PocketBase", err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Alarm deleted successfully"})
	utils.LogMessage("INFO", "Alarm deleted successfully", payload.ID)
}

// ListAlarmsHandler handles listing all alarms
func ListAlarmsHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	alarms, err := pocketbase.ListAlarms(baseURL)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to list alarms"})
		utils.LogMessage("ERROR", "Failed to list alarms from PocketBase", err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, alarms)
	utils.LogMessage("INFO", "Listed all alarms successfully", alarms)
}

// SetAlarmStatusHandler handles updating the activation status of an alarm
func SetAlarmStatusHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		utils.LogMessage("ERROR", "Invalid request method for SetAlarmStatusHandler", nil)
		return
	}

	var payload struct {
		ID        string `json:"id"`
		Activated bool   `json:"activated"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		utils.LogMessage("ERROR", "Failed to decode request body for SetAlarmStatusHandler", err.Error())
		return
	}

	err := pocketbase.SetAlarmStatus(baseURL, payload.ID, payload.Activated)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update alarm status"})
		utils.LogMessage("ERROR", "Failed to update alarm status in PocketBase", err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Alarm status updated successfully"})
	utils.LogMessage("INFO", "Alarm status updated successfully", payload)
}

// ChangeAlarmPlaylistHandler handles changing the playlist of an alarm
func ChangeAlarmPlaylistHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		utils.LogMessage("ERROR", "Invalid request method for ChangeAlarmPlaylistHandler", nil)
		return
	}

	var payload struct {
		ID         string `json:"id"`
		PlaylistID string `json:"playlistId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		utils.LogMessage("ERROR", "Failed to decode request body for ChangeAlarmPlaylistHandler", err.Error())
		return
	}

	err := pocketbase.ChangeAlarmPlaylist(baseURL, payload.ID, payload.PlaylistID)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to change alarm playlist"})
		utils.LogMessage("ERROR", "Failed to change alarm playlist in PocketBase", err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Alarm playlist updated successfully"})
	utils.LogMessage("INFO", "Alarm playlist updated successfully", payload)
}

// ChangeAlarmHourHandler handles changing the hour of an alarm
func ChangeAlarmHourHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, map[string]string{"error": "Invalid request method"})
		utils.LogMessage("ERROR", "Invalid request method for ChangeAlarmHourHandler", nil)
		return
	}

	var payload struct {
		ID   string `json:"id"`
		Hour string `json:"hour"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		utils.WriteJSONResponse(w, http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
		utils.LogMessage("ERROR", "Failed to decode request body for ChangeAlarmHourHandler", err.Error())
		return
	}

	err := pocketbase.ChangeAlarmHour(baseURL, payload.ID, payload.Hour)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, map[string]string{"error": "Failed to change alarm hour"})
		utils.LogMessage("ERROR", "Failed to change alarm hour in PocketBase", err.Error())
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, map[string]string{"message": "Alarm hour updated successfully"})
	utils.LogMessage("INFO", "Alarm hour updated successfully", payload)
}