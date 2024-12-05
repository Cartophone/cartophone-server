package owntone

import (
	"fmt"
	"net/http"
)

// Play sends a request to the Owntone API to start playback
func Play(baseURL string) error {
	url := fmt.Sprintf("%s/api/player/play", baseURL)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to send play command: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}

// Pause sends a request to the Owntone API to pause playback
func Pause(baseURL string) error {
	url := fmt.Sprintf("%s/api/player/pause", baseURL)

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		return fmt.Errorf("failed to send pause command: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response: %s", resp.Status)
	}

	return nil
}