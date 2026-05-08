package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"gorm.io/gorm"
)

type GlobalParams struct {
	logger *slog.Logger
	db     *gorm.DB
}

func HealthHandler(logger *slog.Logger, db *gorm.DB) *GlobalParams {
	return &GlobalParams{logger, db}
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

func (g *GlobalParams) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	dbIsUp := "available"

	if g.db == nil {
		dbIsUp = "unavailable"
	}

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"db":     dbIsUp,
		"time":   time.Now().UTC().Format(time.RFC3339),
	}); err != nil {
		g.logger.Error("Erro na resposta do health!")
	}
}
