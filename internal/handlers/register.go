package handlers

import (
	"fmt"
	"time"
)

// HandleRegisterAction listens for detected card UID and simulates the "Register" action.
func HandleRegisterAction(cardDetectedChan <-chan string) {
	// Simulate waiting for a card to be detected
	for {
		select {
		case uid := <-cardDetectedChan:
			// Simulate registering the card by printing a message to the console
			fmt.Printf("Registering card %s\n", uid)
		}
	}
}