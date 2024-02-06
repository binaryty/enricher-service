package response

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type AgeResponse struct {
	Age uint `json:"age"`
}

type GenderResponse struct {
	Gender string `json:"gender"`
}

type NationalityResponse struct {
	CountryID   string  `json:"country_id"`
	Probability float32 `json:"probability"`
}

type NationalitysResponse struct {
	Countries []NationalityResponse `json:"country"`
}

type IDResponse struct {
	ID int64 `json:"id"`
}

func BadRequest(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusBadRequest, data)
}

func InternalServerError(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusInternalServerError, data)
}

func NotFound(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusNotFound, data)
}

func SuccessfullyCreated(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, data)
}
