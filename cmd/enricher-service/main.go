package main

import (
	"context"
	"github.com/binaryty/enricher-service/internal/app"
	"github.com/binaryty/enricher-service/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

//	@title			Enricher Service Swagger API
//	@version		1.0
//	@description	Swagger API for Golang Project Enricher Service.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Pavel Timochovich
//	@contact.email	t1m0kh0v14@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8082
// @BasePath	/
func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)
	logger.Info("start logging")

	a := app.New(cfg, logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger.Info("starting application", slog.String("host", cfg.HTTPServer.Address))

	go func() {
		a.MustRun(cfg.HTTPServer.Address)
	}()

	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	a.Stop(ctx)

	logger.Info("application stopped")
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
