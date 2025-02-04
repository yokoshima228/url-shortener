package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/yokoshima228/url-shortener/config"
	"github.com/yokoshima228/url-shortener/http-server/handlers/redirect"
	"github.com/yokoshima228/url-shortener/http-server/handlers/rmurl"
	"github.com/yokoshima228/url-shortener/http-server/handlers/url"
	"github.com/yokoshima228/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/yokoshima228/url-shortener/internal/lib/logger/sl"
	"github.com/yokoshima228/url-shortener/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)

	log.Info("Started url-shortener", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Get("/{alias}", redirect.New(log, storage))
	r.Post("/url", url.New(log, storage))
	r.Delete("/url/{alias}", rmurl.New(log, storage))

	log.Info("Started server", slog.String("address", cfg.HttpServer.Address))

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      r,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.Timeout,
		ReadTimeout:  cfg.HttpServer.Timeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		log.Error("Failed to start server", sl.Err(err))
	}

	log.Error("Server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettyLog()
	case envDev:
		log = setupPrettyLog()
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	return log
}

func setupPrettyLog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
