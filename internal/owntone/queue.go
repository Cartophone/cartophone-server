package owntone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cartophone-server/internal/utils"
)

// QueueItem represents an item in the Owntone queue
type QueueItem struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Artist   string `json:"artist"`
	Album    string `json:"album"`
	Duration int    `json:"duration"`
}

// GetQueue retrieves the current Owntone queue
func GetQueue(baseURL string) ([]QueueItem, error) {
	url := fmt.Sprintf("%s/api/player/queue", baseURL)

	utils.LogMessage("INFO", "Fetching Owntone queue", map[string]string{"url": url})

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch queue: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	utils.LogMessage("DEBUG", "Response body received", map[string]string{"body": string(body)})

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status)
	}

	var queue []QueueItem
	if err := json.Unmarshal(body, &queue); err != nil {
		return nil, fmt.Errorf("failed to decode queue response: %w", err)
	}

	return queue, nil
}

// ClearQueue clears the Owntone queue
func ClearQueue(baseURL string) error {
	url := fmt.Sprintf("%s/api/player/queue/clear", baseURL)

	utils.LogMessage("INFO", "Clearing Owntone queue", map[string]string{"url": url})

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create clear queue request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to clear queue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	utils.LogMessage("INFO", "Owntone queue cleared successfully", nil)
	return nil
}

// AddToQueue adds a track to the Owntone queue
func AddToQueue(baseURL string, trackURI string) error {
	url := fmt.Sprintf("%s/api/player/queue/add", baseURL)

	payload := map[string]string{"uri": trackURI}
	data, _ := json.Marshal(payload)

	utils.LogMessage("INFO", "Adding track to Owntone queue", map[string]string{"url": url, "uri": trackURI})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to add track to queue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	utils.LogMessage("INFO", "Track added to Owntone queue successfully", map[string]string{"uri": trackURI})
	return nil
}