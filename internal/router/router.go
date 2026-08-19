package router

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/riazahmedshah/vfs/internal/handler"
	"github.com/riazahmedshah/vfs/internal/middleware"
	"github.com/riazahmedshah/vfs/internal/validation"
)

func NewRouter(h *handler.Handlers) *echo.Echo {
	router := echo.New()

	router.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	router.HTTPErrorHandler = middleware.ErrMiddleware()
	router.Validator = validation.NewCustomValidator()

	router.Group("/api/v1")

	return router
}
