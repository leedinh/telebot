package main

import (
	"log/slog"
	"os"

	config "github.com/leedinh/telebot/bitly/internal/config/url-shortener"
	"github.com/leedinh/telebot/bitly/internal/lib/logger/sl"
	"github.com/leedinh/telebot/bitly/internal/storage/sqlite"
)

func main() {
	// Load the configuration
	cfg := config.LoadConfig()
	log := initLogger(cfg.Env)
	log.Info("Starting the application")
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Failed to create a new storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

}

func initLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "local":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case "dev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	case "prod":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
