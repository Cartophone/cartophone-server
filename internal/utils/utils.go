package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteJSONResponse sends a JSON response to the client and logs the response.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error": "Failed to encode response: %v"}`, err), http.StatusInternalServerError)
		fmt.Printf("[ERROR] Failed to encode response: %v\n", err)
		return
	}

	_, writeErr := w.Write(response)
	if writeErr != nil {
		fmt.Printf("[ERROR] Failed to write response to client: %v\n", writeErr)
		return
	}

	// Log the response payload and status code
	fmt.Printf("[RESPONSE] Status: %d, Payload: %s\n", statusCode, response)
}