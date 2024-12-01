package utils

import (
	"fmt"
	"sync"

	"cartophone-server/internal/constants"
	"cartophone-server/internal/handlers"
)

// StartModeManager manages the application's mode of operation.
func StartModeManager(
	modeSwitch <-chan string,
	cardDetectedChan <-chan string,
	currentMode *string,
	modeLock *sync.Mutex,
	pocketBaseURL string,
) {
	go func() {
		for {
			select {
			case mode := <-modeSwitch:
				fmt.Printf("[DEBUG] modeSwitch signal received: %s\n", mode)

				modeLock.Lock()
				if *currentMode != mode {
					*currentMode = mode
					if mode == constants.ReadMode {
						fmt.Println("[DEBUG] Switched to Read Mode")
					} else if mode == constants.AssociateMode {
						fmt.Println("[DEBUG] Switched to Associate Mode")
					}
				} else {
					fmt.Printf("[DEBUG] Ignoring duplicate signal for mode: %s\n", mode)
				}
				modeLock.Unlock()

			case uid := <-cardDetectedChan:
				modeLock.Lock()
				if *currentMode == constants.ReadMode {
					fmt.Printf("[DEBUG] Detected card in Read Mode: %s\n", uid)
					handlers.HandleReadAction(uid, pocketBaseURL)
				} else {
					fmt.Printf("[DEBUG] Ignoring card %s because we are in Associate Mode\n", uid)
				}
				modeLock.Unlock()
			}
		}
	}()
}