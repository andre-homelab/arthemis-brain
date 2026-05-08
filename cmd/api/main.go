package main

import (
	"log/slog"
	"net/http"
	"os"

	"arthemis-brain/internal/database"
	"arthemis-brain/internal/handlers"
	ownMiddleware "arthemis-brain/internal/middlewares"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
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
		logger.Error("Erro ao inicializar o banco!")
		db = nil
	}

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if _, err := w.Write([]byte("Hello World!")); err != nil {
			logger.Error("Erro na API")
		}
	})

	healthHandler := handlers.HealthHandler(logger, db)
	r.Get("/health", healthHandler.HealthCheck)

	proponentHandler := handlers.ProponentHandler(logger, db)
	r.Post("/proponent/create", proponentHandler.CreateProponent)
	r.Get("/proponent/get", proponentHandler.GetProponent)
	r.Put("/proponent/update", proponentHandler.UpdateProponent)
	r.Delete("/proponent/delete", proponentHandler.DeleteProponent)

	projectHandler := handlers.ProjectHandler(logger, db)
	r.Post("/proponent/create", projectHandler.CreateProject)
	r.Get("/proponent/get", projectHandler.GetProject)
	r.Put("/proponent/update", projectHandler.UpdateProject)
	r.Delete("/proponent/delete", projectHandler.DeleteProject)

	logger.Info("Servidor iniciado!")
	logger.Info("http://localhost:8081")

	if err := http.ListenAndServe(":8081", r); err != nil {
		logger.Error("Erro no roteamento HTTP")
	}
}
