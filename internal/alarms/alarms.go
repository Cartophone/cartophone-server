package alarms

import (
	"fmt"
	"time"

	"cartophone-server/internal/pocketbase"
	"cartophone-server/internal/utils"
)

// StartAlarmChecker starts a goroutine to periodically check for active alarms.
func StartAlarmChecker(baseURL string) {
	go func() {
		for {
			now := time.Now()
			currentTime := fmt.Sprintf("%02d:%02d", now.Hour(), now.Minute())

			// Fetch active alarms for the current time
			alarms, err := pocketbase.FetchActiveAlarms(baseURL, currentTime)
			if err != nil {
				utils.LogMessage("ERROR", "Failed to fetch activated alarms", err.Error())
				time.Sleep(1 * time.Minute) // Retry after a minute
				continue
			}

			// Process each active alarm
			for _, alarm := range alarms {
				playlist, err := pocketbase.GetPlaylist(baseURL, alarm.PlaylistID)
				if err != nil {
					utils.LogMessage("ERROR", fmt.Sprintf("Failed to fetch playlist for alarm %s", alarm.ID), err.Error())
					continue
				}

				// Simulate playing the playlist
				message := fmt.Sprintf("Playing playlist '%s' (URI: %s) for alarm %s", playlist.Name, playlist.URI, alarm.ID)
				utils.LogMessage("ALARM", message, nil)
			}

			// Wait a minute before checking again
			time.Sleep(1 * time.Minute)
		}
	}()
}