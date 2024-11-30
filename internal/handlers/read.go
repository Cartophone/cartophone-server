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

// HandleReadAction listens for detected card UID and simulates the "Play" action.
func HandleReadAction(cardDetectedChan <-chan string, modeSwitch <-chan bool) {
	// Start reading NFC cards
	for {
		select {
		case active := <-modeSwitch:
			// Wait for the mode to switch to true (read mode)
			if active {
				fmt.Println("Reading mode active")
			} else {
				fmt.Println("Reading mode paused")
			}
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