package handlers

import (
	"fmt"
	"cartophone-server/internal/nfc"
)

// ReadHandler handles the /read route
func ReadHandler(cardDetectedChan <-chan string) {
	// Listen for card UID and trigger the action (in the future, connect to Owntone)
	for {
		select {
		case uid := <-cardDetectedChan:
			// Simulate the action of playing the card
			fmt.Printf("Playing card %s\n", uid)
			// Here you can add the logic to trigger Owntone API
		}
	}
}