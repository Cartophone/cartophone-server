package owntone

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"cartophone-server/internal/utils"
)

// QueueItem represents a track in the Owntone queue
type QueueItem struct {
	ID           int    `json:"id"`
	Position     int    `json:"position"`
	Title        string `json:"title"`
	Artist       string `json:"artist"`
	Album        string `json:"album"`
	URI          string `json:"uri"`
	ArtworkURL   string `json:"artwork_url"`
	LengthMillis int    `json:"length_ms"`
}

// FetchQueue fetches the current queue from Owntone
func FetchQueue(baseURL string) ([]map[string]interface{}, error) {
	url := fmt.Sprintf("%s/api/player/queue", baseURL)
	utils.LogMessage("INFO", "Fetching Owntone queue", map[string]string{"url": url})

	resp, err := http.Get(url)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to fetch queue", map[string]string{"error": err.Error()})
		return nil, fmt.Errorf("failed to fetch queue: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body) // Read for debugging purposes
	utils.LogMessage("DEBUG", "Response body received", map[string]string{"body": string(body)})

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response: %s", resp.Status)
	}

	var result struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result.Items, nil
}

// ClearQueue clears the Owntone queue
func ClearQueue(baseURL string) error {
	url := fmt.Sprintf("%s/api/player/queue/clear", baseURL)
	utils.LogMessage("INFO", "Clearing Owntone queue", map[string]string{"url": url})

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to create PUT request to clear queue", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to create clear queue request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to clear queue", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to clear queue: %w", err)
	}
	defer resp.Body.Close()

	// Treat 204 as a successful response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response status", map[string]string{"status": resp.Status, "body": string(body)})
		return fmt.Errorf("unexpected response: %s", string(body))
	}

	utils.LogMessage("INFO", "Owntone queue cleared successfully", nil)
	return nil
}

// AddToQueue adds items to the Owntone queue
func AddToQueue(baseURL string, uris []string) error {
	url := fmt.Sprintf("%s/api/player/queue/items/add", baseURL)
	utils.LogMessage("INFO", "Adding items to Owntone queue", map[string]interface{}{"uris": uris})

	payload := map[string]interface{}{
		"uris": uris,
	}
	data, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		utils.LogMessage("ERROR", "Failed to add track to queue", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to add items to queue: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body) // Read for debugging purposes
	utils.LogMessage("DEBUG", "Response body received", map[string]string{"body": string(body)})

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	utils.LogMessage("INFO", "Owntone queue items added successfully", map[string]interface{}{"uris": uris})
	return nil
}