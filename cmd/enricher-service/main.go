package main

import (
	"github.com/binaryty/enricher-service/internal/app"
	"github.com/binaryty/enricher-service/internal/config"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// init config
	cfg := config.MustLoad()

	// init logger
	logger := setupLogger(cfg.Env)
	logger.Info("start logging")

	// init app
	a := app.New(cfg, logger)

	// start app
	logger.Info("starting application", slog.String("host", cfg.HTTPServer.Address))
	a.MustRun(cfg.HTTPServer.Address)
}

// setupLogger ...
func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
