package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/handler"
	"github.com/riazahmedshah/vfs/internal/middleware"
)

func registerUserRoutes(r *echo.Group, h *handler.Handlers, middleware *middleware.AuthMiddleware) {
	r.GET("/user/me", h.User.GetUser, middleware.RequireAuth())
	auth := r.Group("/auth")
	auth.POST("/google", h.User.GoogleAuth)
	auth.POST("/guest", h.User.GuestAuth)
}
