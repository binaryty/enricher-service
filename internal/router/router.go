package router

import (
	"context"
	"errors"
	"github.com/binaryty/enricher-service/internal/models"
	"github.com/binaryty/enricher-service/internal/response"
	"github.com/binaryty/enricher-service/internal/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
	Echo    *echo.Echo
	service PeopleService
}

// New returns a new instance of Router.
func New(echo *echo.Echo, service PeopleService) *Router {
	return &Router{
		Echo:    echo,
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
		if errors.Is(err, storage.ErrNotFound) {
			return response.NotFound(c, err)
		}
		return response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusOK, person)
}

// Update ...
func (r *Router) Update(c echo.Context) error {
	req := models.Person{}

	if err := c.Bind(&req); err != nil {
		return response.BadRequest(c, err)
	}

	if err := r.service.Update(c.Request().Context(), &req); err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Success(c, response.IDResponse{ID: req.ID})
}

// SelectAll ...
func (r *Router) SelectAll(c echo.Context) error {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		return response.BadRequest(c, err)
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		return response.BadRequest(c, err)
	}

	params := models.Params{
		Limit:  limit,
		Offset: offset,
	}

	persons, err := r.service.SelectAll(c.Request().Context(), params)
	if err != nil {
		return response.InternalServerError(c, err)
	}

	return response.Success(c, persons)
}

// DeleteByID ...
func (r *Router) DeleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.BadRequest(c, err)
	}
	if err := r.service.DeleteByID(c.Request().Context(), int64(id)); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return response.NotFound(c, err)
		}

		return response.InternalServerError(c, err)
	}

	return c.JSON(http.StatusNoContent, nil)
}

func (r *Router) Route() {
	r.Echo.Use(middleware.Recover())

	r.Echo.POST("/people", r.AddPerson)
	r.Echo.GET("people", r.SelectAll)
	r.Echo.GET("/people/:id", r.SelectByID)
	r.Echo.DELETE("/people/:id", r.DeleteByID)
	r.Echo.PUT("/people", r.Update)
}
