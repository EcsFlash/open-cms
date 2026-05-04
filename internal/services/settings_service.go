package services

import (
	"headless-cms/internal/models"
	"headless-cms/internal/repos"
)

type SettingsService struct {
	repo repos.ISettingsRepo
}

type ISettingsService interface {
	Get() (*models.Settings, error)
	Update(s *models.Settings) error
}

func NewSettingsService(repo repos.ISettingsRepo) *SettingsService {
	return &SettingsService{repo: repo}
}

func (s *SettingsService) Get() (*models.Settings, error) {
	return s.repo.Get()
}

func (s *SettingsService) Update(settings *models.Settings) error {
	return s.repo.Update(settings)
}
