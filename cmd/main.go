package main

import (
	"fmt"
	"log"
	"net/http"

	"cartophone-server/internal/api"  // Ensure this points to your internal API package
	"cartophone-server/internal/nfc"  // Ensure this points to your internal NFC package
)

func main() {
	// Initialize NFC reader
	reader, err := nfc.NewReader("pn532_i2c:/dev/i2c-1:0x24")
	if err != nil {
		log.Fatalf("Failed to initialize NFC reader: %v", err)
	}
	defer reader.Close()

	// Set up the HTTP server and routes
	http.HandleFunc("/register-card", api.RegisterCardHandler(reader))  // Register card handler
	http.HandleFunc("/scan-card", api.ScanCardHandler(reader))          // Scan card handler

	// Start the server
	fmt.Println("Starting the server on :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}