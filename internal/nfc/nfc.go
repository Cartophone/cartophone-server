package nfc

import (
	"fmt"
	"time"

	"github.com/clausecker/nfc/v2"
)

// Reader represents the NFC reader instance.
type Reader struct {
	device *nfc.Device // Pointer to the NFC device
}

// NewReader initializes and returns a new NFC Reader.
func NewReader(devicePath string) (*Reader, error) {
	// Open the NFC device
	dev, err := nfc.Open(devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open NFC device: %v", err)
	}

	// Return a new Reader instance
	return &Reader{device: dev}, nil
}

// Close closes the NFC device connection.
func (r *Reader) Close() {
	if r.device != nil {
		r.device.Close()
	}
}

// Scan polls for NFC tags and returns the tag's UID if found.
func (r *Reader) Scan(modulations []nfc.Modulation, attempts int, period time.Duration) (string, error) {
	// Poll for an NFC target
	count, target, err := r.device.InitiatorPollTarget(modulations, attempts, period)
	if err != nil {
		return "", fmt.Errorf("error polling NFC target: %v", err)
	}
	if count == 0 {
		return "", fmt.Errorf("no NFC target detected")
	}

	// Ensure the target is ISO14443a-compatible
	isoTarget, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		return "", fmt.Errorf("unsupported NFC target type")
	}

	// Return the UID as a formatted string
	return fmt.Sprintf("% X", isoTarget.UID), nil
}

// RegisterMode waits for an NFC card to be scanned within the specified timeout.
func (r *Reader) RegisterMode(timeout time.Duration) (string, error) {
	start := time.Now()
	for time.Since(start) < timeout {
		uid, err := r.Scan([]nfc.Modulation{
			{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
		}, 1, 300*time.Millisecond)
		if err == nil && uid != "" {
			return uid, nil
		}
	}
	return "", fmt.Errorf("no NFC card detected in registration mode")
}