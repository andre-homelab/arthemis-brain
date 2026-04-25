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

	db, err := database.ConnectToDatabase()
	if err != nil {
		db = nil
	}
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	healthHandler := handlers.HealthHandler(db)

	r.Get("/health", healthHandler.HealthCheck)

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(textHandler)

	logger.Info("Servidor iniciado!")
	logger.Info("http://localhost:8081")

	http.ListenAndServe(":8081", r)
}
