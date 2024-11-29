package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"cartophone-server/internal/api"
	"cartophone-server/internal/nfc"
	"cartophone-server/internal/owntone"
	"cartophone-server/internal/pocketbase"

	"github.com/go-chi/chi/v5"
)

func main() {
	// Initialize dependencies
	nfcReader, err := nfc.NewReader("pn532_i2c:/dev/i2c-1:0x24")
	if err != nil {
		log.Fatalf("Failed to initialize NFC reader: %v", err)
	}
	defer nfcReader.Close()

	pbClient := pocketbase.NewClient("http://127.0.0.1:8090")
	owntoneClient := owntone.NewClient("http://127.0.0.1:3689")

	// Start background tasks
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		api.RunNFCPoller(nfcReader, pbClient, owntoneClient)
	}()

	go func() {
		defer wg.Done()
		api.RunAlarmMonitor(pbClient, owntoneClient)
	}()

	// Set up HTTP API
	r := chi.NewRouter()
	r.Post("/register-card", func(w http.ResponseWriter, r *http.Request) {
		api.RegisterCardHandler(w, r, nfcReader, pbClient)
	})

	// Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()
	log.Println("HTTP server running on :8080")

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("Shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}
	wg.Wait()
	log.Println("Application exited.")
}
