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
