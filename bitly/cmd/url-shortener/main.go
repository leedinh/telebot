package main

import (
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	config "github.com/leedinh/telebot/bitly/internal/config/url-shortener"
	mwLogger "github.com/leedinh/telebot/bitly/internal/http-server/middleware/logger"
	"github.com/leedinh/telebot/bitly/internal/lib/logger/handlers/slogpretty"
	"github.com/leedinh/telebot/bitly/internal/lib/logger/sl"
	"github.com/leedinh/telebot/bitly/internal/storage/sqlite"
	"golang.org/x/exp/slog"
)

func main() {
	// Load the configuration
	cfg := config.LoadConfig()
	log := initLogger(cfg.Env)
	log.Info(
		"starting url-shortener",
		slog.String("env", cfg.Env),
		slog.String("version", "123"),
	)
	log.Debug("debug messages are enabled")

	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to create a new storage", sl.Err(err))
		os.Exit(1)
	}
	log.Info("Storage has been created")
	_ = storage

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = setupPrettySlog()

	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
