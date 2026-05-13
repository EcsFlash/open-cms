package handler

import (
	"encoding/json"
	"headless-cms/internal/models"
	"headless-cms/pkg/sl"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CreateNews
//
// @Tags Новости
// @Summary Создать новость
// @Description Доступно moderator и admin. Автор новости назначается автоматически из JWT; author_id из тела запроса игнорируется. Краткое описание (short) автоматически вычисляется из markdown-контента (картинки вырезаются regexp).
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.News true "Новость (id игнорируется)"
// @Success 201 {object} models.News
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news [post]
func (h *Handler) CreateNews(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	var n models.News
	if err := json.NewDecoder(c.Request().Body).Decode(&n); err != nil {
		h.log.Error("decode create news", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	if err := h.uc.CreateNews(actor, &n); err != nil {
		h.log.Error("create news failed", sl.Err(err))
		return useCaseErrorResponse(c, err)
	}
	return c.JSON(http.StatusCreated, n)
}

// UpdateNews
//
// @Tags Новости
// @Summary Обновить новость
// @Description Admin может обновить любую новость. Moderator может обновить только свою новость. При обновлении short пересчитывается из content_markdown.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID новости"
// @Param body body models.News true "Новость (id из path)"
// @Success 200 {object} models.News
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news/{id} [patch]
func (h *Handler) UpdateNews(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	var n models.News
	if err := json.NewDecoder(c.Request().Body).Decode(&n); err != nil {
		h.log.Error("decode update news", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	n.ID = uint(id)
	if err := h.uc.UpdateNews(actor, &n); err != nil {
		h.log.Error("update news failed", sl.Err(err))
		return useCaseErrorResponse(c, err)
	}
	updated, _ := h.uc.GetNewsByID(uint(id))
	return c.JSON(http.StatusOK, updated)
}

// DeleteNews
//
// @Tags Новости
// @Summary Удалить новость
// @Description Admin может удалить любую новость. Moderator может удалить только свою новость.
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID новости"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news/{id} [delete]
func (h *Handler) DeleteNews(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	if err := h.uc.DeleteNews(actor, uint(id)); err != nil {
		h.log.Error("delete news failed", sl.Err(err))
		return useCaseErrorResponse(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// GetNews
//
// @Tags Новости
// @Summary Получить новость по ID
// @Produce json
// @Param id path int true "ID новости"
// @Success 200 {object} models.News
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/news/{id} [get]
func (h *Handler) GetNews(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	n, err := h.uc.GetNewsByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
	}
	return c.JSON(http.StatusOK, n)
}

// ListNews
//
// @Tags Новости
// @Summary Получить список новостей
// @Produce json
// @Success 200 {array} models.News
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/news [get]
func (h *Handler) ListNews(c echo.Context) error {
	res, err := h.uc.ListNews()
	if err != nil {
		h.log.Error("list news failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}
