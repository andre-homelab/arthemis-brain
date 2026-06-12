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

	"arthemis-brain/internal/audit"
	"arthemis-brain/internal/database"
	"arthemis-brain/internal/env"
	"arthemis-brain/internal/handlers"

	_ "arthemis-brain/docs"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	// r.Use(ownMiddleware.PermissionMiddleware)

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(textHandler)

	db, err := database.ConnectToDatabase(logger)
	if err != nil {
		logger.Error("Error initializing database!")
	}

	// Configuração da Auditoria (com url do watcher configurável via env)
	watcherURL := env.GetEnv("WATCHER_URL", "http://localhost:8082")
	watcherClient := audit.NewWatcherClient(watcherURL)
	auditMiddleware := audit.AuditMiddleware(watcherClient)

	// Rotas Públicas (Sem Auditoria)
	healthHandler := handlers.HealthHandler(logger, db)
	r.Get("/health", healthHandler.HealthCheck)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("doc.json"),
	))

	// Rotas Auditadas (Com Auditoria)
	r.Group(func(r chi.Router) {
		r.Use(auditMiddleware)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			if _, err := w.Write([]byte("Hello World!")); err != nil {
				logger.Error("Error")
			}
		})

		proponentHandler := handlers.ProponentHandler(logger, db)
		r.Route("/proponent", func(r chi.Router) {
			r.Post("/create", proponentHandler.CreateProponent)
			r.Get("/{id}", proponentHandler.GetProponent)
			r.Get("/", proponentHandler.GetAllProponents)
			r.Patch("/update/{id}", proponentHandler.UpdateProponent)
			r.Delete("/delete/{id}", proponentHandler.DeleteProponent)
		})

		projectHandler := handlers.ProjectHandler(logger, db)
		r.Route("/project", func(r chi.Router) {
			r.Post("/create", projectHandler.CreateProject)
			r.Get("/{id}", projectHandler.GetProject)
			r.Put("/update/{id}", projectHandler.UpdateProject)
			r.Delete("/delete/{id}", projectHandler.DeleteProject)
		})

		userHandler := handlers.UserHandler(logger, db)
		r.Route("/user", func(r chi.Router) {
			r.Post("/create", userHandler.CreateUser)
			r.Get("/{id}", userHandler.GetUser)
			r.Get("/", userHandler.GetAllUsers)
			r.Patch("/update/{id}", userHandler.UpdateUser)
			r.Delete("/delete/{id}", userHandler.DeleteUser)
		})
	})

	logger.Info("Server started!")
	logger.Info("http://localhost:8081")

	if err := http.ListenAndServe(":8081", r); err != nil {
		logger.Error("HTTP routing error", "error", err)
	}
}
