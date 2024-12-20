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

    // Set up HTTP routes for player management
    http.HandleFunc("/player/status", func(w http.ResponseWriter, r *http.Request) {
        handlers.PlayerStatusHandler(config.OwnToneBaseURL, w, r)
    })
    http.HandleFunc("/player/play", func(w http.ResponseWriter, r *http.Request) {
        handlers.PlayHandler(config.OwnToneBaseURL, w, r)
    })
    http.HandleFunc("/player/pause", func(w http.ResponseWriter, r *http.Request) {
        handlers.PauseHandler(config.OwnToneBaseURL, w, r)
    })

    // Queue Management Endpoints
    http.HandleFunc("/player/queue/list", func(w http.ResponseWriter, r *http.Request) {
        handlers.ListQueueHandler(config.OwnToneBaseURL, w, r)
    })

    http.HandleFunc("/player/queue/clear", func(w http.ResponseWriter, r *http.Request) {
        handlers.ClearQueueHandler(config.OwnToneBaseURL, w, r)
    })

    http.HandleFunc("/player/queue/add", func(w http.ResponseWriter, r *http.Request) {
        handlers.AddToQueueHandler(config.OwnToneBaseURL, w, r)
    })

    // Set up HTTP routes for cards management
    http.HandleFunc("/cards/associate", func(w http.ResponseWriter, r *http.Request) {
        handlers.AssociateCardHandler(cardDetectedChan, modeSwitch, config.PocketBaseURL, w, r)
    })

    // Set up HTTP routes for alarm management
    http.HandleFunc("/alarms/create", func(w http.ResponseWriter, r *http.Request) {
        handlers.CreateAlarmHandler(config.PocketBaseURL, w, r)
    })
    http.HandleFunc("/alarms/delete", func(w http.ResponseWriter, r *http.Request) {
        handlers.DeleteAlarmHandler(config.PocketBaseURL, w, r)
    })
    http.HandleFunc("/alarms/list", func(w http.ResponseWriter, r *http.Request) {
        handlers.ListAlarmsHandler(config.PocketBaseURL, w, r)
    })
    http.HandleFunc("/alarms/set-status", func(w http.ResponseWriter, r *http.Request) {
        handlers.SetAlarmStatusHandler(config.PocketBaseURL, w, r)
    })
    http.HandleFunc("/alarms/change-playlist", func(w http.ResponseWriter, r *http.Request) {
        handlers.ChangeAlarmPlaylistHandler(config.PocketBaseURL, w, r)
    })
    http.HandleFunc("/alarms/change-hour", func(w http.ResponseWriter, r *http.Request) {
        handlers.ChangeAlarmHourHandler(config.PocketBaseURL, w, r)
    })

    // Start the HTTP server
    go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()

    // Keep the main thread alive
    select {}
}