package handler

import "headless-cms/internal/models"

// ErrorResponse универсальный формат ошибки.
type ErrorResponse struct {
	Error string `json:"error" example:"описание ошибки"`
}

type AuthRegisterRequest struct {
	Nickname string `json:"nickname" example:"alekc"`
	Password string `json:"password" example:"qwerty123"`
}

type AuthLoginRequest struct {
	Nickname string `json:"nickname" example:"alekc"`
	Password string `json:"password" example:"qwerty123"`
}

type AuthLoginResponse struct {
	Token string      `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  models.User `json:"user"`
}

type MediaListAllResponse struct {
	Images []models.Image `json:"images"`
	Videos []models.Video `json:"videos"`
}

