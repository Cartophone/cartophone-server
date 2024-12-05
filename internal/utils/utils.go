package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// LogMessage logs messages to the console with a specific type and optional payload.
func LogMessage(logType, message string, payload interface{}) {
	if payload != nil {
		encodedPayload, err := json.Marshal(payload)
		if err != nil {
			fmt.Printf("[%s] %s (Failed to encode payload: %v)\n", logType, message, err)
			return
		}
		fmt.Printf("[%s] %s - Payload: %s\n", logType, message, encodedPayload)
	} else {
		fmt.Printf("[%s] %s\n", logType, message)
	}
}

// WriteJSONResponse sends a JSON response to the client and logs the response.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response, err := json.Marshal(payload)
	if err != nil {
		LogMessage("ERROR", "Failed to encode response", err.Error())
		http.Error(w, fmt.Sprintf(`{"error": "Failed to encode response: %v"}`, err), http.StatusInternalServerError)
		return
	}

	if _, writeErr := w.Write(response); writeErr != nil {
		LogMessage("ERROR", "Failed to write response to client", writeErr.Error())
		return
	}

	LogMessage("RESPONSE", fmt.Sprintf("Status: %d", statusCode), payload)
}