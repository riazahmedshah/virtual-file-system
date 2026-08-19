package handler

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(us *service.UserService) *UserHandler {
	return &UserHandler{
		userService: us,
	}
}

type LoginPayload struct {
	Token string `json:"token"`
}

func (h *UserHandler) GoogleAuth(c echo.Context) error {
	var payload LoginPayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("failed to bind request payload", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}
	user, err := h.userService.CreateUser(c.Request().Context(), payload.Token)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"data":    user,
	})
}
