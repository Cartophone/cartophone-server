package nfc

import (
	"fmt"
	"github.com/clausecker/nfc/v2" // Importing the nfc package
	"time"
)

// Reader represents the NFC reader instance.
type Reader struct {
	device *nfc.Device // A pointer to the NFC device
}

// NewReader initializes and returns a new NFC Reader.
func NewReader(devicePath string) (*Reader, error) {
	// Open the NFC device (based on the device path provided)
	dev, err := nfc.Open(devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open NFC device: %v", err)
	}

	// Return a new Reader instance with the pointer to the opened device
	return &Reader{device: &dev}, nil
}

// Close closes the NFC device connection.
func (r *Reader) Close() {
	if r.device != nil {
		r.device.Close()
	}
}

// Scan polls for NFC tags and returns the tag's UID if found.
func (r *Reader) Scan(modulations []nfc.Modulation, attempts int, period time.Duration) (string, error) {
	// Poll for NFC target (tag) with the given modulation and settings
	count, target, err := r.device.InitiatorPollTarget(modulations, attempts, period)
	if err != nil {
		return "", fmt.Errorf("error polling NFC target: %v", err)
	}
	if count == 0 {
		return "", nil // No tag detected
	}

	// Ensure the target is of type ISO14443a
	isoTarget, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		return "", fmt.Errorf("unsupported NFC target type")
	}

	// Return the UID of the detected target as a formatted string
	return fmt.Sprintf("% X", isoTarget.UID), nil
}

// StartPolling starts the NFC polling in a separate goroutine and sends detected card UID via channel.
func (r *Reader) StartPolling(cardDetectedChan chan<- string) {
	// Define modulation types for polling (ISO14443a, 106 kbps)
	modulations := []nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
	}

	// Continuously poll for NFC tags in a separate goroutine
	go func() {
		for {
			// Poll for a target (NFC card)
			uid, err := r.Scan(modulations, 10, 300*time.Millisecond) // 10 attempts, 300ms polling period
			if err != nil {
				fmt.Printf("Error scanning NFC tag: %v\n", err)
				continue
			}

			// If a tag is detected, send its UID to the main thread through the channel
			if uid != "" {
				cardDetectedChan <- uid
			}

			// Wait before polling again
			time.Sleep(200 * time.Millisecond)
		}
	}()
}