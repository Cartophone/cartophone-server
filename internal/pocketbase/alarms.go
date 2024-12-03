package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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
	// Properly encode the filter query
	filter := url.QueryEscape(fmt.Sprintf("hour='%s' && activated=true", currentTime))
	queryURL := fmt.Sprintf("%s/api/collections/alarms/records?filter=%s", baseURL, filter)

	fmt.Printf("[DEBUG] Fetching alarms with URL: %s\n", queryURL)

	resp, err := http.Get(queryURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch alarms: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body) // Read the response body for debugging
	fmt.Printf("[DEBUG] Response body: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var response struct {
		Items []Alarm `json:"items"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
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

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create alarm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var alarm Alarm
	if err := json.NewDecoder(resp.Body).Decode(&alarm); err != nil {
		return nil, fmt.Errorf("failed to decode alarm response: %w", err)
	}

	return &alarm, nil
}

// DeleteAlarm deletes an alarm by ID
func DeleteAlarm(baseURL, id string) error {
	url := fmt.Sprintf("%s/api/collections/alarms/records/%s", baseURL, id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create delete request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to delete alarm: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	return nil
}

// ListAlarms fetches all alarms
func ListAlarms(baseURL string) ([]Alarm, error) {
	url := fmt.Sprintf("%s/api/collections/alarms/records", baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch alarms: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var response struct {
		Items []Alarm `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode alarms response: %w", err)
	}

	return response.Items, nil
}

// SetAlarmStatus updates the status of an alarm
func SetAlarmStatus(baseURL, id string, activated bool) error {
	url := fmt.Sprintf("%s/api/collections/alarms/records/%s", baseURL, id)

	payload := map[string]interface{}{
		"activated": activated,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create patch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update alarm status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	return nil
}

// ChangeAlarmPlaylist updates the playlist ID of an alarm
func ChangeAlarmPlaylist(baseURL, id, playlistID string) error {
	url := fmt.Sprintf("%s/api/collections/alarms/records/%s", baseURL, id)

	payload := map[string]interface{}{
		"playlistId": playlistID,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create patch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update alarm playlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	return nil
}

// ChangeAlarmHour updates the hour of an alarm
func ChangeAlarmHour(baseURL, id, hour string) error {
	url := fmt.Sprintf("%s/api/collections/alarms/records/%s", baseURL, id)

	payload := map[string]interface{}{
		"hour": hour,
	}
	data, _ := json.Marshal(payload)

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to create patch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update alarm hour: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	return nil
}