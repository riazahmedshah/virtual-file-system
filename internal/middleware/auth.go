package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	jwtoken "github.com/riazahmedshah/vfs/internal/model/jwt"
	"github.com/riazahmedshah/vfs/internal/server"
)

type AuthMiddleware struct {
	server *server.Server
}

func NewAuthMiddleware(s *server.Server) *AuthMiddleware {
	return &AuthMiddleware{
		server: s,
	}
}

func (auth *AuthMiddleware) RequireAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cookie, err := c.Cookie("access_token")
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "cookie not found",
				})
			}
			tokenStr := cookie.Value
			claims := &jwtoken.CustomClaims{}
			token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (any, error) {
				return []byte(auth.server.Config.Auth.JwtSecret), nil
			})
			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "invalid or expired token",
				})
			}

			c.Set("userID", claims.UserID)
			c.Set("storageLimit", claims.StorageLimit)
			c.Set("fileSizeLimit", claims.FileSizeLimit)
			return next(c)
		}
	}
}
