package handler

import (
	"errors"
	"headless-cms/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func CurrentActor(c echo.Context) (usecase.Actor, bool) {
	userID, ok := ExtractUserID(c)
	if !ok {
		return usecase.Actor{}, false
	}
	role, ok := ExtractRole(c)
	if !ok {
		return usecase.Actor{}, false
	}
	return usecase.Actor{UserID: userID, Role: role}, true
}

func useCaseErrorResponse(c echo.Context, err error) error {
	switch {
	case errors.Is(err, usecase.ErrForbidden):
		return c.JSON(http.StatusForbidden, map[string]any{"error": "forbidden"})
	case errors.Is(err, gorm.ErrRecordNotFound):
		return c.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
	default:
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
}
