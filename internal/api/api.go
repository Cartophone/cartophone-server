package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cartophone-server/internal/nfc"
	"cartophone-server/internal/owntone"
	"cartophone-server/internal/pocketbase"
)

func RunNFCPoller(reader *nfc.Reader, pbClient *pocketbase.Client, otClient *owntone.Client) {
	modulations := []nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
	}
	for {
		uid, err := reader.Scan(modulations, 10, 300*time.Millisecond)
		if err != nil {
			continue
		}

		playlistID, exists := pbClient.GetPlaylistForUID(uid)
		if exists {
			log.Printf("Playing playlist for UID %s", uid)
			otClient.PlayPlaylist(playlistID)
		}
	}
}

func RunAlarmMonitor(pbClient *pocketbase.Client, otClient *owntone.Client) {
	for {
		alarms := pbClient.GetActiveAlarms()
		now := time.Now()

		for _, alarm := range alarms {
			if now.Hour() == alarm.Hour && now.Minute() == alarm.Minute {
				log.Printf("Triggering alarm: %s", alarm.PlaylistID)
				otClient.PlayPlaylist(alarm.PlaylistID)
			}
		}
		time.Sleep(1 * time.Minute)
	}
}

func RegisterCardHandler(w http.ResponseWriter, r *http.Request, reader *nfc.Reader, pbClient *pocketbase.Client) {
	uid, err := reader.RegisterMode(10 * time.Second)
	if err != nil {
		http.Error(w, "Failed to register card: "+err.Error(), http.StatusRequestTimeout)
		return
	}

	data := map[string]string{"uid": uid}
	json.NewEncoder(w).Encode(data)

	if pbClient.UIDExists(uid) {
		pbClient.UpdateUID(uid)
	} else {
		pbClient.RegisterUID(uid)
	}
}
