package nfc

import (
	"fmt"
	"github.com/clausecker/nfc/v2"  // Correct nfc package import
	"time"
)

// Reader represents the NFC reader instance.
type Reader struct {
	device *nfc.Device // Store a pointer to nfc.Device
}

// NewReader initializes and returns a new NFC Reader.
func NewReader(devicePath string) (*Reader, error) {
	// Open the NFC device using the nfc package
	dev, err := nfc.Open(devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open NFC device: %v", err)
	}

	// Return the Reader struct, passing the pointer to the device
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
	// Poll for NFC target (tag)
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

	// Return the UID of the detected target
	return fmt.Sprintf("% X", isoTarget.UID), nil
}