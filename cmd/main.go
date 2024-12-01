package main

import (
    "fmt"
    "log"
    "net/http"
    "sync"

    "cartophone-server/config"
    "cartophone-server/internal/constants"
    "cartophone-server/internal/handlers"
    "cartophone-server/internal/nfc"
    "cartophone-server/internal/alarms"
)

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
    currentMode := constants.ReadMode

    // Use StartModeManager from handlers
    handlers.StartModeManager(
        modeSwitch,
        cardDetectedChan,
        &currentMode,
        &modeLock,
        config.PocketBaseURL,
    )

    // Start polling for NFC cards
    go reader.StartRead(cardDetectedChan)

    // Start the alarm checker from utils
    alarms.StartAlarmChecker(config.PocketBaseURL)

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