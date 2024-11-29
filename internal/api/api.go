package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"cartophone-server/internal/nfc"  // Ensure correct import of your NFC package
	"cartophone-server/internal/pocketbase"  // Ensure correct import of your Pocketbase package
	"cartophone-server/internal/owntone"    // Ensure correct import of your Owntone package
)

// RegisterCardHandler handles the NFC card registration process
func RegisterCardHandler(w http.ResponseWriter, r *http.Request, reader *nfc.Reader, pbClient *pocketbase.Client) {
	uid, err := reader.RegisterMode(10 * time.Second)
	if err != nil {
		http.Error(w, "Failed to register card: "+err.Error(), http.StatusRequestTimeout)
		return
	}

	// Send the UID as a response
	data := map[string]string{"uid": uid}
	json.NewEncoder(w).Encode(data)

	// Check if the UID exists in Pocketbase
	if pbClient.UIDExists(uid) {
		pbClient.UpdateUID(uid)
	} else {
		pbClient.RegisterUID(uid)
	}
}

// ScanCardHandler scans the NFC card and launches the associated playlist
func ScanCardHandler(w http.ResponseWriter, r *http.Request, reader *nfc.Reader, owntoneClient *owntone.Client) {
	uid, err := reader.Scan([]nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
	}, 10, 300*time.Millisecond)
	if err != nil {
		http.Error(w, "Failed to scan card: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Check if the UID exists in Pocketbase and trigger Owntone playback
	if owntoneClient.IsRegistered(uid) {
		owntoneClient.PlayPlaylist(uid)
		fmt.Fprintf(w, "Playing playlist for UID: %s", uid)
	} else {
		http.Error(w, "Card not registered", http.StatusNotFound)
	}
}