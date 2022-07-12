// Package health provides ...
package controllers

import (
	"encoding/json"
	"net/http"
)

/* Health check if the app stand up */
func Health(w http.ResponseWriter, r *http.Request) {
	response := make(map[string]string)
	response["Message"] = "This is the health controller"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
