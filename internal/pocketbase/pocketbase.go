package pocketbase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Card struct {
	ID       string `json:"id,omitempty"`
	UID      string `json:"uid"`
	Playlist string `json:"playlist,omitempty"`
}

// CheckCard checks if a card exists in the PocketBase database
func CheckCard(baseURL, uid string) (*Card, error) {
	url := fmt.Sprintf("%s/api/collections/cards/records?filter=uid='%s'", baseURL, uid)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to check card: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Card does not exist
	} else if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(body))
	}

	var result struct {
		Items []Card `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Items) == 0 {
		return nil, nil // Card does not exist
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

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	return nil
}