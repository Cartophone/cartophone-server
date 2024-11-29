package handlers

import (
	"fmt"
)

// HandleRegisterAction listens for detected card UID and simulates the "Register" action.
func HandleRegisterAction(cardDetectedChan <-chan string) {
	for {
		select {
		case uid := <-cardDetectedChan:
			// Simulate registering the card by printing a message to the console
			fmt.Printf("Registering card %s\n", uid)
		}
	}
}