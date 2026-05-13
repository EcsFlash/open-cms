package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

// UploadImage
//
// @Tags Медиа
// @Summary Загрузить изображение
// @Description Доступно moderator и admin. Владелец изображения назначается автоматически из JWT. Загружает оригинал + генерирует копии large/medium/mini/thumbnail. Данные хранятся в MinIO и метаданные в Postgres.
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл изображения"
// @Param name formData string false "Имя изображения (alt). Если не задано, берётся имя файла."
// @Success 201 {object} models.Image
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media/images [post]
func (h *Handler) UploadImage(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "file is required"})
	}
	src, err := fh.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	defer src.Close()

	name := c.FormValue("name")
	contentType := fh.Header.Get("Content-Type")

	img, err := h.uc.UploadImage(actor, c.Request().Context(), src, fh.Filename, contentType, name)
	if err != nil {
		return useCaseErrorResponse(c, err)
	}
	return c.JSON(http.StatusCreated, img)
}

// UploadVideo
//
// @Tags Медиа
// @Summary Загрузить видео
// @Description Доступно moderator и admin. Владелец видео назначается автоматически из JWT. Загружает видео в MinIO. Если thumbnail_image_id не задан, thumbnail генерируется автоматически из первого кадра (ffmpeg) и сохраняется как Image того же владельца.
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Файл видео"
// @Param name formData string false "Название видео. Если не задано, берётся имя файла."
// @Param thumbnail_image_id formData int false "ID существующей картинки (Image) как thumbnail"
// @Success 201 {object} models.Video
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media/videos [post]
func (h *Handler) UploadVideo(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	fh, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "file is required"})
	}
	src, err := fh.Open()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	defer src.Close()

	name := c.FormValue("name")
	contentType := fh.Header.Get("Content-Type")

	var thumbID *uint
	if v := strings.TrimSpace(c.FormValue("thumbnail_image_id")); v != "" {
		id64, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid thumbnail_image_id"})
		}
		id := uint(id64)
		thumbID = &id
	}

	video, err := h.uc.UploadVideo(actor, c.Request().Context(), src, fh.Filename, contentType, name, thumbID)
	if err != nil {
		return useCaseErrorResponse(c, err)
	}
	return c.JSON(http.StatusCreated, video)
}

// ListMedia
//
// @Tags Медиа
// @Summary Получить медиа-материалы
// @Description Возвращает изображения, видео или всё сразу.
// @Produce json
// @Param type query string false "Фильтр: images | videos | all" Enums(images,videos,all) default(all)
// @Success 200 {object} MediaListAllResponse "Если type=all"
// @Success 200 {array} models.Image "Если type=images"
// @Success 200 {array} models.Video "Если type=videos"
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/media [get]
func (h *Handler) ListMedia(c echo.Context) error {
	t := strings.ToLower(strings.TrimSpace(c.QueryParam("type")))
	switch t {
	case "", "all":
		images, err := h.uc.ListImages()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
		videos, err := h.uc.ListVideos()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]any{"images": images, "videos": videos})
	case "images":
		images, err := h.uc.ListImages()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, images)
	case "videos":
		videos, err := h.uc.ListVideos()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
		}
		return c.JSON(http.StatusOK, videos)
	default:
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "type must be images, videos, or all"})
	}
}

// GetImage
//
// @Tags Медиа
// @Summary Получить изображение по ID
// @Produce json
// @Param id path int true "ID изображения"
// @Success 200 {object} models.Image
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/images/{id} [get]
func (h *Handler) GetImage(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	img, err := h.uc.GetImageByID(uint(id64))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
	}
	return c.JSON(http.StatusOK, img)
}

// GetVideo
//
// @Tags Медиа
// @Summary Получить видео по ID
// @Produce json
// @Param id path int true "ID видео"
// @Success 200 {object} models.Video
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/videos/{id} [get]
func (h *Handler) GetVideo(c echo.Context) error {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	v, err := h.uc.GetVideoByID(uint(id64))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
	}
	return c.JSON(http.StatusOK, v)
}

// PatchImage
//
// @Tags Медиа
// @Summary Переименовать изображение
// @Description Admin может переименовать любое изображение. Moderator может переименовать только свое изображение. Обновляет поле name (alt/подпись), файлы в MinIO не перезагружаются.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID изображения"
// @Param body body MediaRenameBody true "Новое имя"
// @Success 200 {object} models.Image
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/images/{id} [patch]
func (h *Handler) PatchImage(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	if err := h.uc.RenameImage(actor, uint(id64), body.Name); err != nil {
		return useCaseErrorResponse(c, err)
	}
	img, err := h.uc.GetImageByID(uint(id64))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
	}
	return c.JSON(http.StatusOK, img)
}

// DeleteImageHandler
//
// @Tags Медиа
// @Summary Удалить изображение
// @Description Admin может удалить любое изображение. Moderator может удалить только свое изображение. Удаляет запись и объект(ы) в MinIO, включая производные размеры.
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID изображения"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/images/{id} [delete]
func (h *Handler) DeleteImageHandler(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	if err := h.uc.DeleteImage(actor, c.Request().Context(), uint(id64)); err != nil {
		return useCaseErrorResponse(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// PatchVideo
//
// @Tags Медиа
// @Summary Переименовать видео
// @Description Admin может переименовать любое видео. Moderator может переименовать только свое видео. Обновляет поле name; файл в MinIO не перезагружается.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "ID видео"
// @Param body body MediaRenameBody true "Новое название"
// @Success 200 {object} models.Video
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/videos/{id} [patch]
func (h *Handler) PatchVideo(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	var body struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(c.Request().Body).Decode(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	if err := h.uc.RenameVideo(actor, uint(id64), body.Name); err != nil {
		return useCaseErrorResponse(c, err)
	}
	v, err := h.uc.GetVideoByID(uint(id64))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]any{"error": "not found"})
	}
	return c.JSON(http.StatusOK, v)
}

// DeleteVideoHandler
//
// @Tags Медиа
// @Summary Удалить видео
// @Description Admin может удалить любое видео. Moderator может удалить только свое видео. Удаляет запись и объект в MinIO; thumbnail-изображение не удаляется автоматически.
// @Security BearerAuth
// @Produce json
// @Param id path int true "ID видео"
// @Success 204 {string} string "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/videos/{id} [delete]
func (h *Handler) DeleteVideoHandler(c echo.Context) error {
	actor, ok := CurrentActor(c)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]any{"error": "unauthorized"})
	}
	id64, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]any{"error": "invalid id"})
	}
	if err := h.uc.DeleteVideo(actor, c.Request().Context(), uint(id64)); err != nil {
		return useCaseErrorResponse(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

