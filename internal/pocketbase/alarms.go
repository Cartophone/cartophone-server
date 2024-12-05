package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"cartophone-server/internal/utils"
)

// Alarm represents an alarm object in PocketBase
type Alarm struct {
	ID         string `json:"id"`
	Hour       string `json:"hour"`
	Activated  bool   `json:"activated"`
	PlaylistID string `json:"playlistId"`
}

// FetchActiveAlarms fetches alarms based on the current time and activation status
func FetchActiveAlarms(baseURL, currentTime string) ([]Alarm, error) {
	filter := url.QueryEscape(fmt.Sprintf("hour='%s' && activated=true", currentTime))
	queryURL := fmt.Sprintf("%s/api/collections/alarms/records?filter=%s", baseURL, filter)

	utils.LogMessage("INFO", "Fetching active alarms", nil)

	resp, err := http.Get(queryURL)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to fetch alarms", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("failed to fetch alarms: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	utils.LogMessage("DEBUG", "Response body received", nil)

	if resp.StatusCode != http.StatusOK {
		utils.LogMessage("ERROR", "Unexpected response while fetching alarms", map[string]interface{}{"status": resp.StatusCode, "body": string(body)})
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var response struct {
		Items []Alarm `json:"items"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		utils.LogMessage("ERROR", "Failed to decode alarms response", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("failed to decode alarms response: %w", err)
	}

	return response.Items, nil
}

// CreateAlarm creates a new alarm in PocketBase
func CreateAlarm(baseURL, playlistID, hour string, activated bool) (*Alarm, error) {
	url := fmt.Sprintf("%s/api/collections/alarms/records", baseURL)

	payload := map[string]interface{}{
		"playlistId": playlistID,
		"hour":       hour,
		"activated":  activated,
	}
	data, _ := json.Marshal(payload)

	utils.LogMessage("INFO", "Creating a new alarm", map[string]interface{}{"payload": payload})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		utils.LogMessage("ERROR", "Failed to create alarm", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("failed to create alarm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response while creating alarm", map[string]interface{}{"status": resp.StatusCode, "body": string(body)})
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var alarm Alarm
	if err := json.NewDecoder(resp.Body).Decode(&alarm); err != nil {
		utils.LogMessage("ERROR", "Failed to decode created alarm response", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("failed to decode alarm response: %w", err)
	}

	utils.LogMessage("INFO", "Alarm created successfully", map[string]interface{}{"alarm": alarm})
	return &alarm, nil
}

// DeleteAlarm deletes an alarm by ID
func DeleteAlarm(baseURL, id string) error {
	url := fmt.Sprintf("%s/api/collections/alarms/records/%s", baseURL, id)

	utils.LogMessage("INFO", "Deleting alarm", map[string]interface{}{"id": id})

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to create delete request", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to delete alarm", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to delete alarm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response while deleting alarm", map[string]interface{}{"status": resp.StatusCode, "body": string(body)})
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	utils.LogMessage("INFO", "Alarm deleted successfully", map[string]interface{}{"id": id})
	return nil
}

// ListAlarms fetches all alarms
func ListAlarms(baseURL string) ([]Alarm, error) {
	url := fmt.Sprintf("%s/api/collections/alarms/records", baseURL)

	utils.LogMessage("INFO", "Fetching all alarms", nil)

	resp, err := http.Get(url)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to fetch alarms", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("failed to fetch alarms: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response while listing alarms", map[string]interface{}{"status": resp.StatusCode, "body": string(body)})
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var response struct {
		Items []Alarm `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		utils.LogMessage("ERROR", "Failed to decode alarms response", map[string]interface{}{"error": err.Error()})
		return nil, fmt.Errorf("failed to decode alarms response: %w", err)
	}

	utils.LogMessage("INFO", "Fetched all alarms successfully", map[string]interface{}{"count": len(response.Items)})
	return response.Items, nil
}

// SetAlarmStatus updates the status of an alarm
func SetAlarmStatus(baseURL, id string, activated bool) error {
	url := fmt.Sprintf("%s/api/collections/alarms/records/%s", baseURL, id)

	payload := map[string]interface{}{
		"activated": activated,
	}
	data, _ := json.Marshal(payload)

	utils.LogMessage("INFO", "Updating alarm status", map[string]interface{}{"id": id, "activated": activated})

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	if err != nil {
		utils.LogMessage("ERROR", "Failed to create patch request", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to create patch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to update alarm status", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to update alarm status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response while updating alarm status", map[string]interface{}{"status": resp.StatusCode, "body": string(body)})
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	utils.LogMessage("INFO", "Alarm status updated successfully", map[string]interface{}{"id": id, "activated": activated})
	return nil
}

// ChangeAlarmPlaylist updates the playlist ID of an alarm
func ChangeAlarmPlaylist(baseURL, id, playlistID string) error {
	url := fmt.Sprintf("%s/api/collections/alarms/records/%s", baseURL, id)

	payload := map[string]interface{}{
		"playlistId": playlistID,
	}
	data, _ := json.Marshal(payload)

	utils.LogMessage("INFO", "Changing alarm playlist", map[string]interface{}{"id": id, "playlistId": playlistID})

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	if err != nil {
		utils.LogMessage("ERROR", "Failed to create patch request", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to create patch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to change alarm playlist", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to change alarm playlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response while changing alarm playlist", map[string]interface{}{"status": resp.StatusCode, "body": string(body)})
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	utils.LogMessage("INFO", "Alarm playlist changed successfully", map[string]interface{}{"id": id, "playlistId": playlistID})
	return nil
}

// ChangeAlarmHour updates the hour of an alarm
func ChangeAlarmHour(baseURL, id, hour string) error {
	url := fmt.Sprintf("%s/api/collections/alarms/records/%s", baseURL, id)

	payload := map[string]interface{}{
		"hour": hour,
	}
	data, _ := json.Marshal(payload)

	utils.LogMessage("INFO", "Changing alarm hour", map[string]interface{}{"id": id, "hour": hour})

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	if err != nil {
		utils.LogMessage("ERROR", "Failed to create patch request", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to create patch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to change alarm hour", map[string]interface{}{"error": err.Error()})
		return fmt.Errorf("failed to change alarm hour: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response while changing alarm hour", map[string]interface{}{"status": resp.StatusCode, "body": string(body)})
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	utils.LogMessage("INFO", "Alarm hour changed successfully", map[string]interface{}{"id": id, "hour": hour})
	return nil
}