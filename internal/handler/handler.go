package handler

import (
	"headless-cms/internal/config"
	"headless-cms/internal/usecase"
	"log/slog"
)

type Handler struct {
	cfg *config.Config
	uc  usecase.IUseCase
	log slog.Logger
}

func New(cfg *config.Config, uc usecase.IUseCase, log slog.Logger) *Handler {
	return &Handler{cfg: cfg, uc: uc, log: log}
}

