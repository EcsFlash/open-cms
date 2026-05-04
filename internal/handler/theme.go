package handler

import (
	"encoding/json"
	"headless-cms/internal/models"
	"headless-cms/pkg/sl"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetTheme
//
// @Tags Тема оформления
// @Summary Получить тему (светлая/тёмная палитра)
// @Produce json
// @Success 200 {object} models.Theme
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/theme [get]
func (h *Handler) GetTheme(c echo.Context) error {
	t, err := h.uc.GetTheme()
	if err != nil {
		h.log.Error("get theme", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

// UpdateTheme
//
// @Tags Тема оформления
// @Summary Обновить тему
// @Description Обычно одна запись (id=1). Цвета в формате CSS, например #007aff.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.Theme true "Тема (id из body или существующая запись)"
// @Success 200 {object} models.Theme
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/theme [patch]
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
