package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/handler"
	"github.com/riazahmedshah/vfs/internal/middleware"
)

func registerDirRoutes(r *echo.Group, h *handler.Handlers, middleware *middleware.AuthMiddleware) {
	dir := r.Group("/dir")
	dir.Use(middleware.RequireAuth())
	dir.GET("/:dirId", h.Dir.GetDirectoryContents)
	dir.POST("/:parentId", h.Dir.CreateDirectory)
	dir.PUT("/:dirId", h.Dir.UpdateDirectory)
	dir.DELETE("/:dirId", h.Dir.DeleteDirectory)
}
