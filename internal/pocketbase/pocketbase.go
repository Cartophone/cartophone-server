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
	ID       string `json:"id"`
	UID      string `json:"uid"`
	Playlist string `json:"playlist"` // Field for the associated playlist
}

// CheckCard checks if a card exists in the PocketBase database
func CheckCard(baseURL, uid string) (*Card, error) {
	// Properly encode the filter query
	filter := url.QueryEscape(fmt.Sprintf("uid='%s'", uid))
	url := fmt.Sprintf("%s/api/collections/cards/records?filter=%s", baseURL, filter)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PocketBase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// No card found
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
		// No card found
		return nil, nil
	}

	// Return the first matching card
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