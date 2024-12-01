package handlers

import (
	"fmt"
	"net/http"
)

// WriteResponse writes an HTTP response with a status code and message
func WriteResponse(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(message))
	if err != nil {
		fmt.Printf("Error writing HTTP response: %v\n", err)
	}
}