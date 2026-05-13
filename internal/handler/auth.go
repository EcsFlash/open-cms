package handler

import (
	"encoding/json"
	"errors"
	"headless-cms/internal/models"
	"headless-cms/internal/services"
	"headless-cms/pkg/sl"
	"net/http"

	"github.com/labstack/echo/v4"
)

type registerReq struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type loginReq struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}
type RemoveSupervisorRequest struct {
	Nickname string `json:"nickname"`
}

// Register
//
// @Tags Авторизация
// @Summary Регистрация пользователя
// @Description Создаёт пользователя с ролью user (обычный читатель).
// @Accept json
// @Produce json
// @Param body body AuthRegisterRequest true "Данные регистрации"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register [post]
func (h *Handler) Register(c echo.Context) error {
	var req registerReq
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		h.log.Error("decode register", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	u, err := h.uc.Register(req.Nickname, req.Password)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrNicknameTaken) {
			status = http.StatusConflict
		}
		h.log.Error("register failed", sl.Err(err))
		return c.JSON(status, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, u)
}

// RegisterSupervisor
//
// @Tags Авторизация
// @Summary Защищенный метод для регистрации служебных учетных записей, с предоставлением привелегий роли Moderator.
// @Description Создаёт пользователя с ролью moderator. Доступно только пользователю с ролью admin.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body AuthRegisterRequest true "Данные регистрации служебной учетной записи"
// @Success 201 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register_supervisor [post]
func (h *Handler) RegisterSupervisor(c echo.Context) error {
	var req registerReq
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		h.log.Error("decode register", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	u, err := h.uc.RegisterSupervisor(req.Nickname, req.Password, models.RoleModerator)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrNicknameTaken) {
			status = http.StatusConflict
		}
		h.log.Error("register failed", sl.Err(err))
		return c.JSON(status, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, u)
}

// Login
//
// @Tags Авторизация
// @Summary Вход и получение JWT токена
// @Description Возвращает JWT. Для защищённых методов используйте заголовок Authorization: Bearer <token>.
// @Accept json
// @Produce json
// @Param body body AuthLoginRequest true "Данные входа"
// @Success 200 {object} AuthLoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *Handler) Login(c echo.Context) error {
	var req loginReq
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		h.log.Error("decode login", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	token, u, err := h.uc.Login(req.Nickname, req.Password)
	if err != nil {
		if errors.Is(err, services.ErrInvalidCredentials) {
			return c.JSON(http.StatusUnauthorized, map[string]any{"error": "invalid credentials"})
		}
		h.log.Error("login failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{
		"token": token,
		"user":  u,
	})
}

// RemoveSupervisor
//
// @Tags Авторизация
// @Summary Защищенный метод для удаления служебных учетных записей
// @Description Удаляет пользователя c ролью moderator. Доступно только пользователю с ролью admin.
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body RemoveSupervisorBody true "Данные модератора, подлежащего удалению"
// @Success 200 {object} NoData
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/register_supervisor [post]
func (h *Handler) RemoveSupervisor(c echo.Context) error {
	var req RemoveSupervisorRequest
	if err := json.NewDecoder(c.Request().Body).Decode(&req); err != nil {
		h.log.Error("decode register", sl.Err(err))
		return c.JSON(http.StatusBadRequest, map[string]any{"error": err.Error()})
	}
	err := h.uc.RemoveSupervisor(req.Nickname)
	if err != nil {
		h.log.Error("register failed", sl.Err(err))
		return c.JSON(http.StatusInternalServerError, map[string]any{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, map[string]any{"error": nil})
}
