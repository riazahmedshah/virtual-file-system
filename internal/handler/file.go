package handler

import (
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/model/file"
	"github.com/riazahmedshah/vfs/internal/service"
)

type FileHandler struct {
	fileService *service.FileService
}

func NewFileHandler(fs *service.FileService) *FileHandler {
	return &FileHandler{
		fileService: fs,
	}
}

func (h *FileHandler) UploadAndCreateFile(c echo.Context) error {
	dirID, err := uuid.Parse(c.Param("dirId"))
	if err != nil {
		slog.Error("failed to parse dirId from request params", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid dirId")
	}
	userID := c.Get("userID").(uuid.UUID)
	var payload file.CreateFilePayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("failed to bind request payload", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	if err := c.Validate(payload); err != nil {
		slog.Error("failed to validate request payload", "error", err)
		return err
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		slog.Error("failed to get file from form", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "file is required")
	}

	src, err := fileHeader.Open()
	if err != nil {
		slog.Error("failed to open uploaded file", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to process file")
	}
	defer src.Close()

	fileItem, err := h.fileService.UploadFile(c.Request().Context(), userID, dirID, src, &payload)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, fileItem)
}
