package handlers

import (
	"sync"

	"cartophone-server/internal/constants"
	"cartophone-server/internal/utils"
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
				utils.LogMessage("INFO", "Mode switch signal received", map[string]interface{}{"mode": mode})

				modeLock.Lock()
				if *currentMode != mode {
					*currentMode = mode
					if mode == constants.ReadMode {
						utils.LogMessage("INFO", "Switched to Read Mode", nil)
					} else if mode == constants.AssociateMode {
						utils.LogMessage("INFO", "Switched to Associate Mode", nil)
					}
				} else {
					utils.LogMessage("INFO", "Duplicate mode switch signal ignored", map[string]interface{}{"mode": mode})
				}
				modeLock.Unlock()

			case uid := <-cardDetectedChan:
				modeLock.Lock()
				if *currentMode == constants.ReadMode {
					utils.LogMessage("INFO", "Card detected in Read Mode", map[string]interface{}{"uid": uid})
					HandleReadAction(uid, pocketBaseURL) // Direct call to HandleReadAction
				} else {
					utils.LogMessage("INFO", "Card ignored because of current mode", map[string]interface{}{
						"uid":  uid,
						"mode": *currentMode,
					})
				}
				modeLock.Unlock()
			}
		}
	}()
}