package nfc

import (
	"fmt"
	"github.com/clausecker/nfc/v2"  // Correct nfc package import
	"time"
)

// Reader represents the NFC reader instance.
type Reader struct {
	device *nfc.Device // Pointer to the NFC device
}

// NewReader initializes and returns a new NFC Reader.
func NewReader(devicePath string) (*Reader, error) {
	dev, err := nfc.Open(devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open NFC device: %v", err)
	}

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
	count, target, err := r.device.InitiatorPollTarget(modulations, attempts, period)
	if err != nil {
		return "", fmt.Errorf("error polling NFC target: %v", err)
	}
	if count == 0 {
		return "", nil // No tag detected
	}

	isoTarget, ok := target.(*nfc.ISO14443aTarget)
	if !ok {
		return "", fmt.Errorf("unsupported NFC target type")
	}

	return fmt.Sprintf("% X", isoTarget.UID), nil
}