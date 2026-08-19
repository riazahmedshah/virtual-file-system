package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/handler"
)

func registerUserRoutes(r *echo.Group, h *handler.Handlers) {
	auth := r.Group("/auth")
	auth.POST("/google", h.User.GoogleAuth)
}
