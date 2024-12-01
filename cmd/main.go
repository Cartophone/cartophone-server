package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"

    "cartophone-server/config"
    "cartophone-server/internal/handlers"
    "cartophone-server/internal/nfc"
)

const (
    ReadMode      = "read"
    AssociateMode = "associate"
)

func startModeManager(modeSwitch <-chan string, cardDetectedChan <-chan string, currentMode *string, modeLock *sync.Mutex) {
    go func() {
        for {
            select {
            case mode := <-modeSwitch:
                // Log every mode switch
                fmt.Printf("[DEBUG] modeSwitch signal received: %s\n", mode)

                modeLock.Lock()
                if *currentMode != mode {
                    *currentMode = mode
                    if mode == ReadMode {
                        fmt.Println("[DEBUG] Switched to Read Mode")
                    } else if mode == AssociateMode {
                        fmt.Println("[DEBUG] Switched to Associate Mode")
                    }
                } else {
                    fmt.Printf("[DEBUG] Ignoring duplicate signal for mode: %s\n", mode)
                }
                modeLock.Unlock()

            case uid := <-cardDetectedChan:
                // Handle card detection in read mode
                modeLock.Lock()
                if *currentMode == ReadMode {
                    fmt.Printf("[DEBUG] Detected card in Read Mode: %s\n", uid)
                    handlers.HandleReadAction(uid, config.PocketBaseURL)
                } else {
                    fmt.Printf("[DEBUG] Ignoring card %s because we are in Associate Mode\n", uid)
                }
                modeLock.Unlock()
            }
        }
    }()
}

func main() {
    // Display a nice start message
    fmt.Println("Cartophone server is starting...")
    fmt.Println("Welcome to Cartophone! Ready to scan NFC cards and interact with Owntone and Pocketbase.")
    fmt.Println("Press Ctrl+C to stop the server.")

    // Load the configuration
    config, err := config.LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Initialize the NFC reader
    reader, err := nfc.NewReader(config.DevicePath)
    if err != nil {
        log.Fatalf("Failed to initialize NFC reader: %v", err)
    }
    defer reader.Close()

    // Channels for NFC card detection and mode switching
    cardDetectedChan := make(chan string)
    modeSwitch := make(chan string)

    // Synchronization for mode state
    var modeLock sync.Mutex
    currentMode := ReadMode

    // Start the mode manager
    startModeManager(modeSwitch, cardDetectedChan, &currentMode, &modeLock)

    // Start polling for NFC cards
    go reader.StartRead(cardDetectedChan)

    // Start alarm checker
    handlers.StartAlarmChecker(config.PocketBaseURL)

    // HTTP endpoint for associate mode
    http.HandleFunc("/associate", func(w http.ResponseWriter, r *http.Request) {
        handlers.AssociateHandler(cardDetectedChan, modeSwitch, config.PocketBaseURL, w, r)
    })

    // Start the HTTP server
    go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()

    // Keep the main thread alive
    select {}
}