package repos

import (
	"headless-cms/internal/models"

	"gorm.io/gorm"
)

type SettingsRepo struct {
	db *gorm.DB
}

func NewSettingsRepo(db *gorm.DB) *SettingsRepo {
	return &SettingsRepo{db: db}
}

type ISettingsRepo interface {
	Get() (*models.Settings, error)
	Update(s *models.Settings) error
}

func (r *SettingsRepo) Get() (*models.Settings, error) {
	var s models.Settings
	if err := r.db.First(&s, 1).Error; err != nil {
		def := models.DefaultSettings
		if err := r.db.Create(&def).Error; err != nil {
			return nil, err
		}
		return &def, nil
	}
	return &s, nil
}

func (r *SettingsRepo) Update(s *models.Settings) error {
	return r.db.Model(&models.Settings{}).Where("id = 1").Updates(map[string]any{
		"project_name":     s.ProjectName,
		"logo_url":         s.LogoURL,
		"favicon_url":      s.FaviconURL,
		"meta_title":       s.MetaTitle,
		"meta_description": s.MetaDescription,
		"meta_keywords":    s.MetaKeywords,
		"og_image_url":     s.OgImageURL,
		"banner_text":      s.BannerText,
		"banner_image_url": s.BannerImageURL,
	}).Error
}
