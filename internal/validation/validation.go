package validation

import (
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func NewCustomValidator() *CustomValidator {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name, _, _ := strings.Cut(field.Tag.Get("json"), ",")
		if name == "" {
			return ""
		}
		return name
	})

	return &CustomValidator{
		Validator: v,
	}
}

func (cv *CustomValidator) Validate(i any) error {
	err := cv.Validator.Struct(i)
	if err == nil {
		return nil
	}

	// Check if the error is a collection of validation errors
	if castedErrors, ok := err.(validator.ValidationErrors); ok {
		errorResponse := make(map[string]string)

		for _, fieldErr := range castedErrors {
			fieldName := fieldErr.Field()

			switch fieldErr.Tag() {
			case "required":
				errorResponse[fieldName] = "This field is required"
			case "email":
				errorResponse[fieldName] = "Invalid email format"
			case "max":
				errorResponse[fieldName] = "Value exceeds maximum allowed length of " + fieldErr.Param()
			case "oneof":
				errorResponse[fieldName] = "Must be one of the following: " + fieldErr.Param()
			default:
				errorResponse[fieldName] = "Invalid value (failed " + fieldErr.Tag() + " restriction)"
			}
		}

		return echo.NewHTTPError(http.StatusBadRequest, map[string]any{
			"message": "validation failed",
			"errors":  errorResponse,
		})
	}

	return echo.NewHTTPError(http.StatusInternalServerError, "internal server error")
}
