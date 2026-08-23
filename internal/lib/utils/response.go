package utils

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func SuccessResponse(c echo.Context, status int, message string, data any) error {
	return c.JSON(status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c echo.Context, status int, message string) error {
	return c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}

func ValidationErrorResponse(c echo.Context, message string, fieldErrors map[string]string) error {
	return c.JSON(http.StatusBadRequest, Response{
		Success: false,
		Message: message,
		Data:    fieldErrors,
	})
}
