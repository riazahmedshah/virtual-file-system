package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/handler"
)

func Registerv1Routes(router *echo.Group, h *handler.Handlers) {
	registerUserRoutes(router, h)
}
