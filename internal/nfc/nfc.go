package nfc

import (
    "fmt"
    "github.com/clausecker/nfc/v2" // Correct import for nfc
)

type Reader struct {
    device *nfc.Device
}

func NewReader(devicePath string) (*Reader, error) {
    dev, err := nfc.Open(devicePath)
    if err != nil {
        return nil, fmt.Errorf("failed to open NFC device: %v", err)
    }
    return &Reader{device: dev}, nil
}

func (r *Reader) Scan(modulations []nfc.Modulation, attempts int, period time.Duration) (string, error) {
    count, target, err := r.device.InitiatorPollTarget(modulations, attempts, period)
    if err != nil || count == 0 {
        return "", fmt.Errorf("no NFC target detected")
    }

    isoTarget, ok := target.(*nfc.ISO14443aTarget)
    if !ok {
        return "", fmt.Errorf("unsupported NFC target type")
    }

    return fmt.Sprintf("% X", isoTarget.UID), nil
}