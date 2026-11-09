package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// writeJSON sends a JSON response with an HTTP status code.
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Println("Error writing JSON response:", err)
	}
}

// serverError logs the real error and sends a generic 500 response,
// so database details are never shown to the client.
func serverError(w http.ResponseWriter, err error) {
	log.Println("Internal error:", err)
	writeJSON(w, http.StatusInternalServerError, map[string]string{
		"error": "Internal server error",
	})
}
