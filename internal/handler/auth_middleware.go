package handler

import (
	"headless-cms/internal/config"
	"headless-cms/internal/models"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	ctxUserRoleKey = "user_role"
	ctxUserIDKey   = "user_id"
)

func JWTAuth(cfg *config.Config) echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(cfg.JWT.Secret),
	})
}

func ExtractRole(c echo.Context) (models.Role, bool) {
	tok, ok := c.Get("user").(*jwt.Token)
	if !ok || tok == nil {
		return "", false
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return "", false
	}
	roleAny, ok := claims["role"]
	if !ok {
		return "", false
	}
	roleStr, ok := roleAny.(string)
	if !ok || roleStr == "" {
		return "", false
	}
	return models.Role(roleStr), true
}

func ExtractUserID(c echo.Context) (uint, bool) {
	tok, ok := c.Get("user").(*jwt.Token)
	if !ok || tok == nil {
		return 0, false
	}
	claims, ok := tok.Claims.(jwt.MapClaims)
	if !ok {
		return 0, false
	}
	subAny, ok := claims["sub"]
	if !ok {
		return 0, false
	}
	// MapClaims numbers are float64 after JSON decoding.
	switch v := subAny.(type) {
	case float64:
		return uint(v), true
	default:
		return 0, false
	}
}

func RequireRole(allowed ...models.Role) echo.MiddlewareFunc {
	allowedSet := make(map[string]bool, len(allowed))
	for _, r := range allowed {
		allowedSet[string(r)] = true
	}
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// echo-jwt stores token under "user". We'll parse role on demand in handler later
			role, ok := ExtractRole(c)
			if !ok {
				return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
			}
			if !allowedSet[string(role)] {
				return c.JSON(http.StatusForbidden, map[string]any{"error": "forbidden"})
			}
			return next(c)
		}
	}
}

