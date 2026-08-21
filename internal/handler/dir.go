package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/model/dir"
	"github.com/riazahmedshah/vfs/internal/service"
)

type DirHandler struct {
	dirService *service.DirService
}

func NewDirHandler(ds *service.DirService) *DirHandler {
	return &DirHandler{
		dirService: ds,
	}
}

func (h *DirHandler) CreateDirectory(c echo.Context) error {
	parentID, err := uuid.Parse(c.Param("parentId"))
	if err != nil {
		slog.Error("failed to parse parentId from request params", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid parentId")
	}
	userID := c.Get("userID").(uuid.UUID)
	var payload dir.CreateDirPayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("failed to bind request payload", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(payload); err != nil {
		slog.Error("failed to validate request payload", "error", err)
		return err
	}

	directory, err := h.dirService.CreateDirectory(c.Request().Context(), userID, parentID, &payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, directory)
}

func (h *DirHandler) GetDirectoryContents(c echo.Context) error {
	dirID, err := uuid.Parse(c.Param("parentId"))
	if err != nil {
		slog.Error("failed to parse dirId from request params", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid dirId")
	}
	userID := c.Get("userID").(uuid.UUID)

	directoryContents, err := h.dirService.GetDirectoryContent(c.Request().Context(), userID, dirID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, directoryContents)
}

func (h *DirHandler) UpdateDirectory(c echo.Context) error {
	dirID, err := uuid.Parse(c.Param("dirId"))
	if err != nil {
		slog.Error("failed to parse dirId from request params", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid dirId")
	}
	userID := c.Get("userID").(uuid.UUID)

	var payload dir.UpdateDirpayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("failed to bind request payload", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(payload); err != nil {
		slog.Error("failed to validate request payload", "error", err)
		return err
	}

	updatedDir, err := h.dirService.UpdateDirectory(c.Request().Context(), userID, dirID, &payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedDir)
}

func (h *DirHandler) DeleteDirectory(c echo.Context) error {
	dirID, err := uuid.Parse(c.Param("dirId"))
	if err != nil {
		slog.Error("failed to parse dirId from request params", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid dirId")
	}
	userID := c.Get("userID").(uuid.UUID)

	err = h.dirService.DeleteDirectory(c.Request().Context(), userID, dirID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "directory deleted successfully"})
}
