package router

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/riazahmedshah/vfs/internal/handler"
	"github.com/riazahmedshah/vfs/internal/middleware"
	v1 "github.com/riazahmedshah/vfs/internal/router/v1"
	"github.com/riazahmedshah/vfs/internal/server"
	"github.com/riazahmedshah/vfs/internal/validation"
)

func NewRouter(s *server.Server, h *handler.Handlers) *echo.Echo {
	router := echo.New()

	router.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PATCH, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders:     []string{echo.HeaderContentType, echo.HeaderAuthorization},
		AllowCredentials: true,
	}))

	authMiddleware := middleware.NewAuthMiddleware(s)
	router.Validator = validation.NewCustomValidator()
	router.HTTPErrorHandler = middleware.ErrMiddleware()

	v1Group := router.Group("/api/v1")
	v1.Registerv1Routes(v1Group, h, authMiddleware)
	return router
}
