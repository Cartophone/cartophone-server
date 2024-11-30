package main

import (
    "cartophone-server/config"
    "cartophone-server/internal/constants" // Import the constants package
    "cartophone-server/internal/handlers"
    "cartophone-server/internal/nfc"
    "fmt"
    "log"
    "net/http"
    "sync"
)

func main() {
    // Display a nice start message
    fmt.Println("Cartophone server is starting...")
    fmt.Println("Welcome to Cartophone! Ready to scan NFC cards and interact with Owntone and Pocketbase.")
    fmt.Println("Press Ctrl+C to stop the server.")

    // Load the configuration from config.json
    config, err := config.LoadConfig("config.json")
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }

    // Initialize the NFC reader using the device path from the config
    reader, err := nfc.NewReader(config.DevicePath)
    if err != nil {
        log.Fatalf("Failed to initialize NFC reader: %v", err)
    }
    defer reader.Close()

    // Create channels for NFC card detection and mode switching
    cardDetectedChan := make(chan string)
    modeSwitch := make(chan string)

    // Synchronization for mode state
    var modeLock sync.Mutex
    currentMode := constants.ReadMode

    // Goroutine to manage NFC card detection and mode state
    go func() {
        for {
            select {
            case mode := <-modeSwitch:
                // Log every mode switch
                fmt.Printf("[DEBUG] modeSwitch signal received: %s\n", mode)

                modeLock.Lock()
                currentMode = mode
                modeLock.Unlock()

                if mode == constants.ReadMode {
                    fmt.Println("Switched to Read Mode")
                } else if mode == constants.AssociateMode {
                    fmt.Println("Switched to Associate Mode")
                }

            case uid := <-cardDetectedChan:
                // Handle card detection in read mode
                modeLock.Lock()
                if currentMode == constants.ReadMode {
                    handlers.HandleReadAction(uid, "http://127.0.0.1:8090")
                } else {
                    fmt.Printf("Ignoring card %s because we are in Associate Mode\n", uid)
                }
                modeLock.Unlock()
            }
        }
    }()

    // Start polling for NFC cards
    go reader.StartRead(cardDetectedChan)

    // HTTP endpoint for associate mode
    http.HandleFunc("/associate", func(w http.ResponseWriter, r *http.Request) {
        modeLock.Lock()
        if currentMode == constants.AssociateMode {
            fmt.Println("Already in associate mode")
            w.WriteHeader(http.StatusConflict)
            fmt.Fprintf(w, "Already in associate mode")
            modeLock.Unlock()
            return
        }

        currentMode = constants.AssociateMode
        modeLock.Unlock()

        // Trigger associate mode
        modeSwitch <- constants.AssociateMode

        // Handle the association logic
        handlers.AssociateHandler(cardDetectedChan, modeSwitch, "http://127.0.0.1:8090", w, r)

        // Revert back to read mode after association is complete
        modeSwitch <- constants.ReadMode
    })

    // Start the HTTP server
    go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()

    // Main thread remains idle
    select {}
}