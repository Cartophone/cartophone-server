package owntone

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetPlayerStatus fetches the status of the OwnTone player
func GetPlayerStatus(baseURL string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/player", baseURL)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch player status: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response from OwnTone: %s", string(body))
	}

	var playerStatus map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&playerStatus); err != nil {
		return nil, fmt.Errorf("failed to decode player status response: %w", err)
	}

	return playerStatus, nil
}