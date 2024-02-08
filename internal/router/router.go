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

	echoSwagger "github.com/swaggo/echo-swagger"
	_ "github.com/swaggo/echo-swagger/example/docs"
)

const (
	PageSize = 10

	StatusOk         = "Ok"
	StatusCrated     = "Successfully Created"
	StatusInternal   = "Internal Error"
	StatusBadRequest = "Bad Request"
	StatusNotFound   = "Not Found"
	StatusNoContent  = "No Content"
)

type PeopleService interface {
	AddPerson(ctx context.Context, rawData models.RawPerson) (int64, error)
	SelectByID(ctx context.Context, id int64) (*models.Person, error)
	Update(ctx context.Context, params *models.Person) error
	SelectAll(ctx context.Context, params models.Params) ([]models.Person, error)
	DeleteByID(ctx context.Context, id int64) error
}

type Router struct {
	service PeopleService
}

// New returns a new instance of Router.
func New(service PeopleService) *Router {
	return &Router{
		service: service,
	}
}

// AddPerson godoc
//
//	@Summary		Add Person
//	@Tags			person
//	@Description	get NSP to enrich it and add
//	@ID				add-person
//	@Accept			json
//	@Produce		json
//	@Param			RawPerson	body		models.RawPerson	true	"name, surname, patronymic"
//	@Success		201			{object}	response.Response
//	@Failure		400			{object}	response.Response
//	@Router			/person [post]
func (r *Router) AddPerson(c echo.Context) error {
	req := models.RawPerson{}

	if err := c.Bind(&req); err != nil {
		return response.SendResponse(c, http.StatusBadRequest, StatusBadRequest, err)
	}

	id, err := r.service.AddPerson(c.Request().Context(), req)

	if err != nil {
		return response.SendResponse(c, http.StatusInternalServerError, StatusInternal, err)
	}

	return response.SendResponse(c, http.StatusCreated, StatusCrated, response.IDResponse{
		ID: id,
	})
}

// SelectByID godoc
//
//	@Summary		Get person by id from storage
//	@Tags			person
//	@Description	get id from url params and find person
//	@ID				get-person
//	@Produce		json
//	@Param			id	path		int	true	"Person ID"
//	@Success		200	{object}	models.Person
//	@Failure		400	{object}	response.Response
//	@Router			/person/{id} [get]
func (r *Router) SelectByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.SendResponse(c, http.StatusBadRequest, StatusBadRequest, err)
	}

	person, err := r.service.SelectByID(c.Request().Context(), int64(id))
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return response.SendResponse(c, http.StatusNotFound, StatusNotFound, err)
		}
		return response.SendResponse(c, http.StatusInternalServerError, StatusInternal, err)
	}

	return response.SendResponse(c, http.StatusOK, StatusOk, person)
}

// Update godoc
//
//	@Summary		Update person in storage
//	@Tags			person
//	@Description	update person
//	@ID				update-person
//	@Accept			json
//	@Produce		json
//	@Param			input	body		models.Person	true	"id, name, surname, patronymic, age, gender, nationality"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Router			/person [put]
func (r *Router) Update(c echo.Context) error {
	req := models.Person{}

	if err := c.Bind(&req); err != nil {
		return response.SendResponse(c, http.StatusBadRequest, StatusBadRequest, err)
	}

	if err := r.service.Update(c.Request().Context(), &req); err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return response.SendResponse(c, http.StatusNotFound, StatusNotFound, err)
		}
		return response.SendResponse(c, http.StatusInternalServerError, StatusInternal, err)
	}

	return response.SendResponse(c, http.StatusOK, StatusOk, response.IDResponse{
		ID: req.ID,
	})
}

// SelectAll godoc
//
//	@Summary		Get a list of persons by params
//	@Tags			person
//	@Description	Get a list of persons based on query parameters
//	@ID				get-all-persons
//	@Accept			json
//	@Produce		json
//	@Param			page	query		int	true	"Id of page of results"
//	@Success		200		{object}	response.Response
//	@Failure		400		{object}	response.Response
//	@Router			/persons [get]
func (r *Router) SelectAll(c echo.Context) error {
	pageId, err := strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return response.SendResponse(c, http.StatusBadRequest, StatusBadRequest, err)
	}

	params := models.Params{
		Limit:  PageSize,
		Offset: (pageId - 1) * PageSize,
	}

	persons, err := r.service.SelectAll(c.Request().Context(), params)
	if err != nil {
		return response.SendResponse(c, http.StatusInternalServerError, StatusInternal, err)
	}

	return response.SendResponse(c, http.StatusOK, StatusOk, persons)
}

// DeleteByID godoc
//
//	@Summary		delete person from storage by id
//	@Tags			person
//	@Description	delete person
//	@ID				delete-person
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Person ID"
//	@Success		204	{object}	response.Response
//	@Failure		400	{object}	response.Response
//	@Router			/person/{id} [delete]
func (r *Router) DeleteByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return response.SendResponse(c, http.StatusBadRequest, StatusBadRequest, err)
	}
	if err := r.service.DeleteByID(c.Request().Context(), int64(id)); err != nil {

		return response.SendResponse(c, http.StatusInternalServerError, StatusInternal, err)
	}

	return response.SendResponse(c, http.StatusNoContent, StatusNoContent, nil)
}

// Route setup router.
func (r *Router) Route(e *echo.Echo) {
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.POST("/person", r.AddPerson)
	e.GET("/persons", r.SelectAll)
	e.GET("/person/:id", r.SelectByID)
	e.DELETE("/person/:id", r.DeleteByID)
	e.PUT("/person", r.Update)
}
