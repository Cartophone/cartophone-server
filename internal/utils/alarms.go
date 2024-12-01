package utils

import (
	"fmt"
	"time"

	"cartophone-server/internal/pocketbase"
)

// StartAlarmChecker starts a goroutine to periodically check for active alarms.
func StartAlarmChecker(baseURL string) {
	go func() {
		for {
			now := time.Now()
			currentTime := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())

			alarms, err := pocketbase.FetchActiveAlarms(baseURL, currentTime)
			if err != nil {
				fmt.Printf("[ERROR] Failed to fetch activated alarms: %v\n", err)
				time.Sleep(1 * time.Minute) // Retry after a minute
				continue
			}

			for _, alarm := range alarms {
				playlist, err := pocketbase.GetPlaylist(baseURL, alarm.PlaylistID)
				if err != nil {
					fmt.Printf("[ERROR] Failed to fetch playlist for alarm %s: %v\n", alarm.ID, err)
					continue
				}

				// Simulate playing the playlist
				fmt.Printf("[ALARM] Playing playlist '%s' (URI: %s) for alarm %s\n", playlist.Name, playlist.URI, alarm.ID)
			}

			time.Sleep(1 * time.Minute) // Wait a minute before checking again
		}
	}()
}