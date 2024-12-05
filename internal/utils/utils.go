package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// WriteJSONResponse sends a JSON response with the provided status code and data.
func WriteJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		fmt.Printf("[DEBUG] Failed to write JSON response: %v\n", err)
	}
}