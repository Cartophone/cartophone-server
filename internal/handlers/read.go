package handlers

import (
	"fmt"
)

// HandleReadAction simulates handling detected cards in read mode.
func HandleReadAction(uid string) {
	fmt.Printf("Playing card %s\n", uid)
}