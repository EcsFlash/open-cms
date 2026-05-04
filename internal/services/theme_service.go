package services

import (
	"headless-cms/internal/models"
	"headless-cms/internal/repos"
)

type ThemeService struct {
	repo repos.IThemeRepo
}

type IThemeService interface {
	Get() (*models.Theme, error)
	Update(t *models.Theme) error
}

func NewThemeService(repo repos.IThemeRepo) *ThemeService {
	return &ThemeService{repo: repo}
}

func (s *ThemeService) Get() (*models.Theme, error) {
	return s.repo.Get()
}

func (s *ThemeService) Update(t *models.Theme) error {
	return s.repo.Update(t)
}
