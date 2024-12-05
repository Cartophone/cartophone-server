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
func FetchQueue(baseURL string) ([]QueueItem, error) {
	url := fmt.Sprintf("%s/api/queue", baseURL)
	utils.LogMessage("INFO", "Fetching Owntone queue", map[string]string{"url": url})

	resp, err := http.Get(url)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to fetch queue", map[string]string{"error": err.Error()})
		return nil, fmt.Errorf("failed to fetch queue: %w", err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	utils.LogMessage("DEBUG", "Response body received", map[string]string{"body": string(body)})

	if resp.StatusCode != http.StatusOK {
		utils.LogMessage("ERROR", "Unexpected response status", map[string]string{"status": resp.Status})
		return nil, fmt.Errorf("unexpected response: %s", resp.Status)
	}

	var response struct {
		Items []QueueItem `json:"items"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		utils.LogMessage("ERROR", "Failed to decode queue response", map[string]string{"error": err.Error()})
		return nil, fmt.Errorf("failed to decode queue response: %w", err)
	}

	return response.Items, nil
}

// ClearQueue clears the current Owntone queue
func ClearQueue(baseURL string) error {
	url := fmt.Sprintf("%s/api/queue/clear", baseURL)
	utils.LogMessage("INFO", "Clearing Owntone queue", map[string]string{"url": url})

	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to create PUT request to clear queue", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to create PUT request to clear queue: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		utils.LogMessage("ERROR", "Failed to clear queue", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to clear queue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response status", map[string]string{"status": resp.Status, "body": string(body)})
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}

// AddToQueue adds a track to the Owntone queue
func AddToQueue(baseURL, uri string) error {
	url := fmt.Sprintf("%s/api/queue/items/add", baseURL)
	payload := map[string]string{"uri": uri}
	data, _ := json.Marshal(payload)

	utils.LogMessage("INFO", "Adding track to Owntone queue", map[string]string{"url": url, "uri": uri})

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		utils.LogMessage("ERROR", "Failed to add track to queue", map[string]string{"error": err.Error()})
		return fmt.Errorf("failed to add track to queue: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		utils.LogMessage("ERROR", "Unexpected response status", map[string]string{"status": resp.Status, "body": string(body)})
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}