package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Card represents a card record in the PocketBase database
type Card struct {
	ID         string `json:"id"`
	UID        string `json:"uid"`
	PlaylistID string `json:"playlistId"` // Links to the Playlist ID
}

// Playlist represents a playlist record in the PocketBase database
type Playlist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URI  string `json:"uri"` // URI to be played
}

// Alarm represents an alarm in the PocketBase "alarms" collection
type Alarm struct {
	ID         string `json:"id"`
	Hour       string `json:"hour"`
	Activated  bool   `json:"activated"`
	PlaylistID string `json:"playlistId"`
}

// CheckCard checks if a card exists in the PocketBase database
func CheckCard(baseURL, uid string) (*Card, error) {
	filter := url.QueryEscape(fmt.Sprintf("uid='%s'", uid))
	url := fmt.Sprintf("%s/api/collections/cards/records?filter=%s", baseURL, filter)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PocketBase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status)
	}

	var result struct {
		Items []Card `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil
	}

	return &result.Items[0], nil
}

// AddCard adds a new card to the PocketBase database
func AddCard(baseURL string, card Card) error {
	url := fmt.Sprintf("%s/api/collections/cards/records", baseURL)

	payload, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("failed to marshal card: %w", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to add card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	return nil
}

// UpdateCard updates an existing card in PocketBase
func UpdateCard(baseURL string, card Card) error {
	url := fmt.Sprintf("%s/api/collections/cards/records/%s", baseURL, card.ID)

	payload, err := json.Marshal(card)
	if err != nil {
		return fmt.Errorf("failed to marshal card: %w", err)
	}

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create update request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to update card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}

// GetPlaylist fetches a playlist by ID from the PocketBase database
func GetPlaylist(baseURL, playlistID string) (*Playlist, error) {
	url := fmt.Sprintf("%s/api/collections/playlists/records/%s", baseURL, playlistID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlist: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var playlist Playlist
	if err := json.NewDecoder(resp.Body).Decode(&playlist); err != nil {
		return nil, fmt.Errorf("failed to decode playlist response: %w", err)
	}

	return &playlist, nil
}

// FetchActiveAlarms fetches all active alarms for a specific time.
func FetchActiveAlarms(baseURL, currentTime string) ([]Alarm, error) {
	url := fmt.Sprintf(
		"%s/api/collections/alarms/records?filter=hour='%s' && activated=true",
		baseURL, currentTime,
	)

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