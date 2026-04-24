package main

import (
	"log/slog"
	"net/http"
	"os"

	"arthemis-brain/internal/handlers"
	ownMiddleware "arthemis-brain/internal/middlewares"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(chiMiddleware.Logger)
	r.Use(ownMiddleware.PermissionMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	r.Get("/health", handlers.HealthCheck)

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(textHandler)

	logger.Info("Servidor iniciado!")
	logger.Info("http://localhost:8081")

	http.ListenAndServe(":8081", r)
}
