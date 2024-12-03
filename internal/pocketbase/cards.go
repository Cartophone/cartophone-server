package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Card represents a card object in PocketBase
type Card struct {
	ID         string `json:"id"`
	UID        string `json:"uid"`
	PlaylistID string `json:"playlistId"`
}

// CheckCard checks if a card exists in PocketBase by UID
func CheckCard(baseURL, uid string) (*Card, error) {
	url := fmt.Sprintf("%s/api/collections/cards/records?filter=uid='%s'", baseURL, uid)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var response struct {
		Items []Card `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode card response: %w", err)
	}

	if len(response.Items) == 0 {
		return nil, nil // Card does not exist
	}

	return &response.Items[0], nil // Return the first matching card
}

// AddCard adds a new card to PocketBase
func AddCard(baseURL, uid, playlistID string) (*Card, error) {
	url := fmt.Sprintf("%s/api/collections/cards/records", baseURL)

	payload := map[string]interface{}{
		"uid":        uid,
		"playlistId": playlistID,
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("failed to add card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var card Card
	if err := json.NewDecoder(resp.Body).Decode(&card); err != nil {
		return nil, fmt.Errorf("failed to decode card response: %w", err)
	}

	return &card, nil
}