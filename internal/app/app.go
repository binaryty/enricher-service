package app

import (
	"context"
	"errors"
	"github.com/binaryty/enricher-service/internal/config"
	"github.com/binaryty/enricher-service/internal/router"
	"github.com/binaryty/enricher-service/internal/services/enricher"
	services "github.com/binaryty/enricher-service/internal/services/people"
	storage "github.com/binaryty/enricher-service/internal/storage/people/postgres"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"os"
)

type App struct {
	l *slog.Logger
	e *echo.Echo
	r *router.Router
}

func New(cfg *config.Config, log *slog.Logger) *App {
	app := App{}

	repo, err := storage.New(cfg)
	if err != nil {
		log.Error("can't initialize storage", slog.String("[ERROR]", err.Error()))
		os.Exit(1)
	}

	enricherSrv := enricher.New(cfg)

	service := services.New(log, repo, enricherSrv)

	app.r = router.New(service)

	app.e = echo.New()

	app.r.Route(app.e)

	return &app
}

func (a *App) MustRun(host string) {
	if err := a.e.Start(host); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.l.Error("can't start application", slog.String("[FATAL]", err.Error()))
		panic(err)
	}
}

func (a *App) Stop(ctx context.Context) {
	if err := a.e.Shutdown(ctx); err != nil {
		a.l.Error("can't stop application", slog.String("[FATAL]", err.Error()))
		panic(err)
	}
}
