package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	//"internal/env"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	connectionString := "host=arthemis-brain-postgres port=5432 user=app_user password=app_password dbname=app_db sslmode=disable"

	_, err := gorm.Open(
		postgres.Open(connectionString),
		&gorm.Config{},
	)

	dbIsUp := "available"
	if err != nil {
		dbIsUp = "unavailable"
	}

	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"db":     dbIsUp,
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}
