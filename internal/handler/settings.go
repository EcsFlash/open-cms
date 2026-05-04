package handler

import (
	"encoding/json"
	"headless-cms/internal/models"
	"headless-cms/pkg/sl"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h *Handler) GetSettings(c echo.Context) error {
	s, err := h.uc.GetSettings()
	if err != nil {
		h.log.Error("get settings", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, s)
}

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
