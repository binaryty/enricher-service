package router

import (
	"context"
	"errors"
	"github.com/binaryty/enricher-service/internal/models"
	"github.com/binaryty/enricher-service/internal/response"
	"github.com/binaryty/enricher-service/internal/storage/people/postgres"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type PeopleService interface {
	AddPerson(ctx context.Context, rawData models.RawPerson) (int64, error)
	SelectByID(ctx context.Context, id int64) (*models.Person, error)
	Update(ctx context.Context, params *models.Person) error
	SelectAll(ctx context.Context, params models.Params) ([]models.Person, error)
	DeleteByID(ctx context.Context, id int64) error
}

type Router struct {
	echo    *echo.Echo
	service PeopleService
}

// New returns a new instance of Router.
func New(echo *echo.Echo, service PeopleService) *Router {
	return &Router{
		echo:    echo,
		service: service,
	}
}

// AddPerson ...
func (r *Router) AddPerson(c echo.Context) error {
	req := models.RawPerson{}

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err)
	}

	id, err := r.service.AddPerson(c.Request().Context(), req)

	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.SuccessfullyCreated(c, response.IDResponse{
		ID: id,
	})
}

// SelectByID ...
func (r *Router) SelectByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, err)
	}

	person, err := r.service.SelectByID(c.Request().Context(), int64(id))
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			return response.NotFound(c, err)
		}
		return response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusOK, person)
}

// Update ...
func (r *Router) Update(ctx context.Context, params *models.Person) error {
	// TODO: implement me
	panic("implement me")
}

// SelectAll ...
func (r *Router) SelectAll(ctx context.Context, params models.Params) ([]models.Person, error) {
	// TODO: implement me
	panic("implement me")
}

// DeleteByID ...
func (r *Router) DeleteByID(ctx context.Context, id int64) error {
	// TODO: implement me
	panic("implement me")
}
