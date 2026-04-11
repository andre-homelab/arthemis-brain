package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
    "status": "ok",
  	"time":   time.Now().UTC().Format(time.RFC3339),
  })
}



