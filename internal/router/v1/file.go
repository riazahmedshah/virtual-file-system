package v1

import (
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/handler"
	"github.com/riazahmedshah/vfs/internal/middleware"
)

func registerFileRoutes(r *echo.Group, h *handler.Handlers, middleware *middleware.AuthMiddleware) {
	file := r.Group("/file")
	file.Use(middleware.RequireAuth())
	file.POST("/upload/:parentId", h.File.UploadAndCreateFile)
}
