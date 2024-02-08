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

type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func BadRequest(c echo.Context, status string, data interface{}) error {
	return c.JSON(http.StatusBadRequest, Response{
		Code:   http.StatusBadRequest,
		Status: status,
		Data:   data,
	})
}

func InternalServerError(c echo.Context, status string, data interface{}) error {
	return c.JSON(http.StatusInternalServerError, Response{
		Code:   http.StatusInternalServerError,
		Status: status,
		Data:   data,
	})
}

func NotFound(c echo.Context, status string, data interface{}) error {
	return c.JSON(http.StatusNotFound, Response{
		Code:   http.StatusNotFound,
		Status: status,
		Data:   data,
	})
}

func SuccessfullyCreated(c echo.Context, status string, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Code:   http.StatusCreated,
		Status: status,
		Data:   data,
	})
}

func Success(c echo.Context, status string, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:   http.StatusOK,
		Status: status,
		Data:   data,
	})
}
