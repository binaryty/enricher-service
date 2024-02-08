package response

import (
	"github.com/labstack/echo/v4"
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

func SendResponse(c echo.Context, code int, status string, data interface{}) error {
	return c.JSON(code, Response{
		Code:   code,
		Status: status,
		Data:   data,
	})
}
