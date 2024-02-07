package app

import (
	"github.com/binaryty/enricher-service/internal/config"
	"github.com/binaryty/enricher-service/internal/router"
	enricher "github.com/binaryty/enricher-service/internal/services/enricher"
	services "github.com/binaryty/enricher-service/internal/services/people"
	storage "github.com/binaryty/enricher-service/internal/storage/people/postgres"
	"github.com/labstack/echo/v4"
	"log/slog"
	"os"
)

type App struct {
	log     *slog.Logger
	service *services.Service
	router  *router.Router
}

func New(cfg *config.Config, log *slog.Logger) *App {
	a := App{}

	repo, err := storage.New(cfg)
	if err != nil {
		log.Error("can't initialize storage", slog.String("[ERROR]", err.Error()))
		os.Exit(1)
	}

	enr := enricher.New(cfg)
	a.service = services.New(log, repo, enr)
	a.router = router.New(echo.New(), a.service)

	a.router.Route()

	return &a
}

func (a *App) MustRun(host string) {
	if err := a.router.Echo.Start(host); err != nil {
		a.log.Error("can't start application", slog.String("[FATAL]", err.Error()))
		panic(err)
	}
}
