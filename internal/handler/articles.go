package handler

import (
	"encoding/json"
	"headless-cms/internal/models"
	"headless-cms/pkg/sl"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// CreateArticle
//
// @Tags Статьи
// @Summary Создать статью
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body models.Article true "Статья (id игнорируется)"
// @Success 201 {object} models.Article
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles [post]
func (h *Handler) CreateArticle(c echo.Context) error {
	var a models.Article
	if err := json.NewDecoder(c.Request().Body).Decode(&a); err != nil {
		h.log.Error("decode create article", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	a.ID = 0
	if err := h.uc.CreateArticle(&a); err != nil {
		h.log.Error("create article failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, a)
}

// UpdateArticle
//
// @Tags Статьи
// @Summary Обновить статью
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID статьи"
// @Param body body models.Article true "Статья (id из path)"
// @Success 200 {object} models.Article
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles/{id} [patch]
func (h *Handler) UpdateArticle(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	var a models.Article
	if err := json.NewDecoder(c.Request().Body).Decode(&a); err != nil {
		h.log.Error("decode update article", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	a.ID = uint(id)
	if err := h.uc.UpdateArticle(&a); err != nil {
		h.log.Error("update article failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	updated, _ := h.uc.GetArticleByID(uint(id))
	return c.JSON(http.StatusOK, updated)
}

// DeleteArticle
//
// @Tags Статьи
// @Summary Удалить статью
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID статьи"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles/{id} [delete]
func (h *Handler) DeleteArticle(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	if err := h.uc.DeleteArticle(uint(id)); err != nil {
		h.log.Error("delete article failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}

// GetArticle
//
// @Tags Статьи
// @Summary Получить статью по ID
// @Produce json
// @Param id path int true "ID статьи"
// @Success 200 {object} models.Article
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/articles/{id} [get]
func (h *Handler) GetArticle(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	a, err := h.uc.GetArticleByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
	}
	return c.JSON(http.StatusOK, a)
}

// ListArticles
//
// @Tags Статьи
// @Summary Получить все статьи
// @Produce json
// @Success 200 {array} models.Article
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/articles [get]
func (h *Handler) ListArticles(c echo.Context) error {
	res, err := h.uc.ListArticles()
	if err != nil {
		h.log.Error("list articles failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

// ListArticlesBySection
//
// @Tags Статьи
// @Summary Получить статьи раздела
// @Produce json
// @Param id path int true "ID раздела"
// @Success 200 {array} models.Article
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/sections/{id}/articles [get]
func (h *Handler) ListArticlesBySection(c echo.Context) error {
	sectionID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	res, err := h.uc.ListArticlesBySection(uint(sectionID))
	if err != nil {
		h.log.Error("list articles by section failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

