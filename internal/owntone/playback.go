package owntone

import (
	"fmt"
	"net/http"
)

// Play sends a request to the Owntone API to start playback
func Play(baseURL string) error {
	url := fmt.Sprintf("%s/api/player/play", baseURL)

	req, err := http.NewRequest(http.MethodPut, url, nil) // Change to PUT method
	if err != nil {
		return fmt.Errorf("failed to create play command request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send play command: %w", err)
	}
	defer resp.Body.Close()

	// Treat 204 No Content as a successful response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}

// Pause sends a request to the Owntone API to pause playback
func Pause(baseURL string) error {
	url := fmt.Sprintf("%s/api/player/pause", baseURL)

	req, err := http.NewRequest(http.MethodPut, url, nil) // Change to PUT method
	if err != nil {
		return fmt.Errorf("failed to create pause command request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send pause command: %w", err)
	}
	defer resp.Body.Close()

	// Treat 204 No Content as a successful response
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}