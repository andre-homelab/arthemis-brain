package main

import (
	"log/slog"
	"os"
  "net/http"
	"arthemis-brain/internal/handlers"
  "github.com/go-chi/chi/v5"
  "github.com/go-chi/chi/v5/middleware"
)

func main() {
  r := chi.NewRouter()
  r.Use(middleware.Logger)
  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
      w.Write([]byte("Hello World!"))
  })

	r.Get("/health", handlers.HealthCheck)

	textHandler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	logger := slog.New(textHandler)

	logger.Info("Servidor iniciado!")
	logger.Info("http://localhost:3000")

	http.ListenAndServe(":3000", r)
}

