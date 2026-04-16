package handler

import (
	"encoding/json"
	"headless-cms/internal/models"
	"headless-cms/internal/services"
	"headless-cms/pkg/sl"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CreateSection
//
// @Tags Разделы
// @Summary Создать раздел
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.Section true "Раздел (id игнорируется)"
// @Success 201 {object} models.Section
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sections [post]
func (h *Handler) CreateSection(c echo.Context) error {
	var s models.Section
	if err := json.NewDecoder(c.Request().Body).Decode(&s); err != nil {
		h.log.Error("decode create section", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	s.ID = 0
	if err := h.uc.CreateSection(&s); err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrSectionCycle {
			status = http.StatusBadRequest
		}
		h.log.Error("create section failed", sl.Err(err))
		return c.JSON(status, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, s)
}

// UpdateSection
//
// @Tags Разделы
// @Summary Обновить раздел
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID раздела"
// @Param body body models.Section true "Раздел (id из path)"
// @Success 200 {object} models.Section
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sections/{id} [patch]
func (h *Handler) UpdateSection(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	var s models.Section
	if err := json.NewDecoder(c.Request().Body).Decode(&s); err != nil {
		h.log.Error("decode update section", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	s.ID = uint(id)
	if err := h.uc.UpdateSection(&s); err != nil {
		status := http.StatusInternalServerError
		if err == services.ErrSectionCycle {
			status = http.StatusBadRequest
		}
		h.log.Error("update section failed", sl.Err(err))
		return c.JSON(status, map[string]any{"error": err.Error()})
	}
	updated, _ := h.uc.GetSectionByID(uint(id))
	return c.JSON(http.StatusOK, updated)
}

// DeleteSection
//
// @Tags Разделы
// @Summary Удалить раздел
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID раздела"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sections/{id} [delete]
func (h *Handler) DeleteSection(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	if err := h.uc.DeleteSection(uint(id)); err != nil {
		h.log.Error("delete section failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// GetSection
//
// @Tags Разделы
// @Summary Получить раздел по ID
// @Produce json
// @Param id path int true "ID раздела"
// @Success 200 {object} models.Section
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/sections/{id} [get]
func (h *Handler) GetSection(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	s, err := h.uc.GetSectionByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
	}
	return c.JSON(http.StatusOK, s)
}

// ListSections
//
// @Tags Разделы
// @Summary Получить все разделы
// @Produce json
// @Success 200 {array} models.Section
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sections [get]
func (h *Handler) ListSections(c echo.Context) error {
	res, err := h.uc.ListSections()
	if err != nil {
		h.log.Error("list sections failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// ListSectionChildren
//
// @Tags Разделы
// @Summary Получить дочерние разделы
// @Produce json
// @Param id path int true "ID родительского раздела"
// @Success 200 {array} models.Section
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sections/{id}/children [get]
func (h *Handler) ListSectionChildren(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	res, err := h.uc.ListSectionChildren(uint(id))
	if err != nil {
		h.log.Error("list section children failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

