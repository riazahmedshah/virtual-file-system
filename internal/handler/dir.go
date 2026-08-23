package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/lib/utils"
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
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid parent directory id")
	}
	userID := c.Get("userID").(uuid.UUID)
	var payload dir.CreateDirPayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("failed to bind request payload", "error", err)
		return err
	}
	if err := c.Validate(payload); err != nil {
		slog.Error("failed to validate request payload", "error", err)
		return err
	}

	directory, err := h.dirService.CreateDirectory(c.Request().Context(), userID, parentID, &payload)
	if err != nil {
		return err
	}
	return utils.SuccessResponse(c, http.StatusCreated, "directory created successfully", directory)
}

func (h *DirHandler) GetDirectoryContents(c echo.Context) error {
	dirID, err := uuid.Parse(c.Param("dirId"))
	if err != nil {
		slog.Error("failed to parse dirId from request params", "error", err)
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid directory id")
	}
	userID := c.Get("userID").(uuid.UUID)

	directoryContents, err := h.dirService.GetDirectoryContent(c.Request().Context(), userID, dirID)
	if err != nil {
		return err
	}
	return utils.SuccessResponse(c, http.StatusOK, "directory contents retrieved successfully", directoryContents)
}

func (h *DirHandler) UpdateDirectory(c echo.Context) error {
	dirID, err := uuid.Parse(c.Param("dirId"))
	if err != nil {
		slog.Error("failed to parse dirId from request params", "error", err)
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid directory id")
	}
	userID := c.Get("userID").(uuid.UUID)

	var payload dir.UpdateDirpayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("failed to bind request payload", "error", err)
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(payload); err != nil {
		slog.Error("failed to validate request payload", "error", err)
		return err
	}

	updatedDir, err := h.dirService.UpdateDirectory(c.Request().Context(), userID, dirID, &payload)
	if err != nil {
		return err
	}
	return utils.SuccessResponse(c, http.StatusOK, "directory updated successfully", updatedDir)
}

func (h *DirHandler) DeleteDirectory(c echo.Context) error {
	dirID, err := uuid.Parse(c.Param("dirId"))
	if err != nil {
		slog.Error("failed to parse dirId from request params", "error", err)
		return utils.ErrorResponse(c, http.StatusBadRequest, "invalid directory id")
	}
	userID := c.Get("userID").(uuid.UUID)

	err = h.dirService.DeleteDirectory(c.Request().Context(), userID, dirID)
	if err != nil {
		return err
	}
	return utils.SuccessResponse(c, http.StatusOK, "directory deleted successfully", nil)
}
