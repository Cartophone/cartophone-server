package handlers

import (
	"fmt"
)

// HandleReadAction listens for detected card UID and simulates the "Play" action.
func HandleReadAction(cardDetectedChan <-chan string, modeSwitch <-chan bool) {
	// Start reading NFC cards
	for {
		select {
		case <-modeSwitch:
			// Wait for the mode to switch to true (read mode)
			fmt.Println("Reading mode active")
		case uid := <-cardDetectedChan:
			// Only read if in active mode
			select {
			case active := <-modeSwitch:
				if active {
					// Simulate playing the card by printing a message to the console
					fmt.Printf("Playing card %s\n", uid)
				}
			default:
				// Do nothing if the mode is not active
			}
		}
	}
}