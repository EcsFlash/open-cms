package handler

import (
	"encoding/json"
	"headless-cms/internal/models"
	"headless-cms/pkg/sl"
	"net/http"

	"github.com/labstack/echo/v4"
)

// GetSettings
//
// @Tags Настройки сайта
// @Summary Получить настройки (название проекта, SEO, баннер)
// @Produce json
// @Success 200 {object} models.Settings
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/settings [get]
func (h *Handler) GetSettings(c echo.Context) error {
	s, err := h.uc.GetSettings()
	if err != nil {
		h.log.Error("get settings", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

// UpdateSettings
//
// @Tags Настройки сайта
// @Summary Обновить настройки
// @Description Обычно одна запись (id=1). URL логотипа, favicon, OG-image и т.д. — строки или null.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.Settings true "Настройки"
// @Success 200 {object} models.Settings
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/settings [patch]
func (h *Handler) UpdateSettings(c echo.Context) error {
	var s models.Settings
	if err := json.NewDecoder(c.Request().Body).Decode(&s); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	if err := h.uc.UpdateSettings(&s); err != nil {
		h.log.Error("update settings", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	updated, err := h.uc.GetSettings()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, updated)
}
