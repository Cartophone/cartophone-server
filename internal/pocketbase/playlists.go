package pocketbase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Playlist represents a playlist object in PocketBase
type Playlist struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	URI  string `json:"uri"`
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