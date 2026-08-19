package middleware

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/riazahmedshah/vfs/internal/errs"
)

func ErrMiddleware() echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		if appErr, ok := errors.AsType[*errs.AppErr](err); ok {
			if appErr.Err != nil {
				slog.Error("internal system error",
					"code", appErr.Status,
					"user_msg", appErr.Message,
					"internal_err", appErr.Err,
					"path", c.Path(),
					"method", c.Request().Method,
				)
			} else {
				slog.Warn("service error",
					"code", appErr.Status,
					"message", appErr.Message,
					"path", c.Path(),
				)
			}

			_ = c.JSON(appErr.Status, map[string]any{
				"success": false,
				"error":   appErr.Message,
			})
			return
		}

		if echoErr, ok := errors.AsType[*echo.HTTPError](err); ok {
			slog.Warn("echo HTTP error", "code", echoErr.Code, "msg", echoErr.Message)

			_ = c.JSON(echoErr.Code, map[string]any{
				"success": false,
				"error":   fmt.Sprintf("%v", echoErr.Message),
			})
			return
		}

		slog.Error("unhandled critical server error", "error", err, "path", c.Path())

		_ = c.JSON(http.StatusInternalServerError, map[string]any{
			"success": false,
			"error":   "An unexpected server error occurred. Please try again later.",
		})
	}
}
