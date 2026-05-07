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

// @title           Health Cheack
// @version         1.0
// @description     Evaluates the status of the server and database
// @termsOfService  http://swagger.io/terms/

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8081
// @BasePath  /health

// @securityDefinitions.basic  BasicAuth

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
