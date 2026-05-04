package handler

import (
	"encoding/json"
	"headless-cms/internal/models"
	"headless-cms/pkg/sl"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetTheme(c echo.Context) error {
	t, err := h.uc.GetTheme()
	if err != nil {
		h.log.Error("get theme", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (h *Handler) UpdateTheme(c echo.Context) error {
	var t models.Theme
	if err := json.NewDecoder(c.Request().Body).Decode(&t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	if err := h.uc.UpdateTheme(&t); err != nil {
		h.log.Error("update theme", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	updated, err := h.uc.GetTheme()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updated)
}
