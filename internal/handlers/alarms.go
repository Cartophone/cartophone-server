package handlers

import (
	"fmt"
	"time"

	"cartophone-server/internal/pocketbase"
)

// StartAlarmChecker checks for activated alarms every minute and triggers associated actions
func StartAlarmChecker(baseURL string) {
	go func() {
		for {
			now := time.Now().Format("15:04") // Current time in HH:mm format
			fmt.Printf("[DEBUG] Checking alarms for %s\n", now)

			// Fetch active alarms
			alarms, err := pocketbase.FetchActiveAlarms(baseURL)
			if err != nil {
				fmt.Printf("[ERROR] Failed to fetch alarms: %v\n", err)
				time.Sleep(1 * time.Minute) // Retry after 1 minute
				continue
			}

			// Check for matching alarms
			for _, alarm := range alarms {
				if alarm.Hour == now {
					fmt.Printf("[DEBUG] Alarm triggered! Hour: %s, Playlist ID: %s\n", alarm.Hour, alarm.PlaylistID)

					// Fetch the playlist associated with the alarm
					playlist, err := pocketbase.GetPlaylist(baseURL, alarm.PlaylistID)
					if err != nil {
						fmt.Printf("[ERROR] Failed to fetch playlist for alarm: %v\n", err)
						continue
					}

					// Simulate playing the playlist
					fmt.Printf("[INFO] Playing playlist: %s (URI: %s)\n", playlist.Name, playlist.URI)
				}
			}

			time.Sleep(1 * time.Minute) // Check again after 1 minute
		}
	}()
}