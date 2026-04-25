package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type DBStore struct {
	db *gorm.DB
}

func HealthHandler(db *gorm.DB) *DBStore {
	return &DBStore{db}
}

func (db *DBStore) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	dbIsUp := "available"

	if db == nil {
		dbIsUp = "unavailable"
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"db":     dbIsUp,
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
