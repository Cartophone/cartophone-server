package handlers

import (
	"fmt"
)

// HandleReadAction listens for detected card UID and simulates the "Play" action.
func HandleReadAction(cardDetectedChan <-chan string) {
	for {
		select {
		case uid := <-cardDetectedChan:
			// Simulate playing the card by printing a message to the console
			fmt.Printf("Playing card %s\n", uid)
		}
	}
}