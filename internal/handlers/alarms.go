package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"cartophone-server/internal/pocketbase"
)

// CreateAlarmHandler creates a new alarm
func CreateAlarmHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		PlaylistID string `json:"playlistId"`
		Hour       string `json:"hour"`
		Activated  bool   `json:"activated"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	alarm, err := pocketbase.CreateAlarm(baseURL, payload.PlaylistID, payload.Hour, payload.Activated)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create alarm: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(alarm)
}

// DeleteAlarmHandler deletes an alarm by ID
func DeleteAlarmHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing alarm ID", http.StatusBadRequest)
		return
	}

	err := pocketbase.DeleteAlarm(baseURL, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to delete alarm: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, "Alarm deleted successfully")
}

// ListAlarmsHandler lists all alarms
func ListAlarmsHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	alarms, err := pocketbase.ListAlarms(baseURL)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to list alarms: %v", err), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(alarms)
}

// SetAlarmStatusHandler updates the status of an alarm
func SetAlarmStatusHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		ID        string `json:"id"`
		Activated bool   `json:"activated"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := pocketbase.SetAlarmStatus(baseURL, payload.ID, payload.Activated)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update alarm status: %v", err), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "Alarm status updated successfully")
}

// ChangeAlarmPlaylistHandler changes the playlist of an alarm
func ChangeAlarmPlaylistHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		ID         string `json:"id"`
		PlaylistID string `json:"playlistId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := pocketbase.ChangeAlarmPlaylist(baseURL, payload.ID, payload.PlaylistID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to change alarm playlist: %v", err), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "Alarm playlist updated successfully")
}

// ChangeAlarmHourHandler changes the hour of an alarm
func ChangeAlarmHourHandler(baseURL string, w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload struct {
		ID   string `json:"id"`
		Hour string `json:"hour"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := pocketbase.ChangeAlarmHour(baseURL, payload.ID, payload.Hour)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to change alarm hour: %v", err), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, "Alarm hour updated successfully")
}