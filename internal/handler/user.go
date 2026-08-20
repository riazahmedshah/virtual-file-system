package handler

import (
	"log/slog"
	"net/http"
	"time"

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
	Token string `json:"token" validate:"required"`
}

func (h *UserHandler) GoogleAuth(c echo.Context) error {
	var payload LoginPayload
	if err := c.Bind(&payload); err != nil {
		slog.Error("failed to bind request payload", "error", err)
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request payload")
	}

	if err := c.Validate(payload); err != nil {
		return err
	}
	userToken, err := h.userService.CreateUser(c.Request().Context(), payload.Token)
	if err != nil {
		return err
	}
	cookie := new(http.Cookie)
	cookie.Name = "access_token"
	cookie.Value = userToken
	cookie.Expires = time.Now().Add(time.Hour * 24)
	cookie.HttpOnly = true
	cookie.Secure = false // Ensures cookie is only sent over HTTPS (Set to false ONLY in local dev if not using HTTPS)
	cookie.SameSite = http.SameSiteLaxMode
	cookie.Path = "/"
	c.SetCookie(cookie)
	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"token":   userToken,
	})
}
