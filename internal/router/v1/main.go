package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/handler"
	"github.com/riazahmedshah/vfs/internal/middleware"
)

func Registerv1Routes(router *echo.Group, h *handler.Handlers, middleware *middleware.AuthMiddleware) {
	registerUserRoutes(router, h)
	registerFileRoutes(router, h, middleware)
}
