// @title           Arthemis Brain API
// @version         1.0
// @description     API service for Arthemis Brain.
// @termsOfService  http://swagger.io/terms/

// @license.name  MIT
// @license.url   http://opensource.org/licenses/MIT

// @host      localhost:8081
// @BasePath  /
package main

import (
	"log/slog"
	"net/http"
	"os"

	"arthemis-brain/internal/database"
	"arthemis-brain/internal/handlers"
	ownMiddleware "arthemis-brain/internal/middlewares"

	_ "arthemis-brain/docs"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(ownMiddleware.PermissionMiddleware)

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(textHandler)

	db, err := database.ConnectToDatabase(logger)
	if err != nil {
		logger.Error("Error initializing database!")
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello World!")); err != nil {
			logger.Error("Error")
		}
	})

	healthHandler := handlers.HealthHandler(logger, db)
	r.Get("/health", healthHandler.HealthCheck)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	proponentHandler := handlers.ProponentHandler(logger, db)
	r.Route("/proponent", func(r chi.Router) {
		r.Post("/{id}", proponentHandler.CreateProponent)
		r.Get("/{id}", proponentHandler.GetProponent)
		r.Patch("/{id}", proponentHandler.UpdateProponent)
		r.Delete("/{id}", proponentHandler.DeleteProponent)
	})

	projectHandler := handlers.ProjectHandler(logger, db)
	r.Route("/project", func(r chi.Router) {
		r.Post("/{id}", projectHandler.CreateProject)
		r.Get("/{id}", projectHandler.GetProject)
		r.Put("/{id}", projectHandler.UpdateProject)
		r.Delete("/{id}", projectHandler.DeleteProject)
	})

	logger.Info("Server started!")
	logger.Info("http://localhost:8081")

	if err := http.ListenAndServe(":8081", r); err != nil {
		logger.Error("HTTP routing error", "error", err)
	}
}
