package nfc

import (
	"fmt"
	"github.com/clausecker/nfc/v2"
	"time"
)

// Reader struct for NFC reader
type Reader struct {
	device *nfc.Device
}

// NewReader initializes the NFC reader.
func NewReader(devicePath string) (*Reader, error) {
	dev, err := nfc.Open(devicePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open NFC device: %v", err)
	}
	return &Reader{device: &dev}, nil
}

// Close closes the NFC device connection
func (r *Reader) Close() {
	if r.device != nil {
		r.device.Close()
	}
}

// StartRead starts scanning NFC tags
func (r *Reader) StartRead(cardDetectedChan chan<- string) {
	modulations := []nfc.Modulation{
		{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106},
	}

	go func() {
		for {
			count, target, err := r.device.InitiatorPollTarget(modulations, 10, 300*time.Millisecond)
			if err != nil {
				fmt.Printf("Error scanning NFC tag: %v\n", err)
				continue
			}
			if count > 0 {
				isoTarget, ok := target.(*nfc.ISO14443aTarget)
				if ok {
					cardDetectedChan <- fmt.Sprintf("% X", isoTarget.UID)
				}
			}
			time.Sleep(1 * time.Second)
		}
	}()
}