package repos

import (
	"headless-cms/internal/models"

	"gorm.io/gorm"
)

type ThemeRepo struct {
	db *gorm.DB
}

func NewThemeRepo(db *gorm.DB) *ThemeRepo {
	return &ThemeRepo{db: db}
}

type IThemeRepo interface {
	Get() (*models.Theme, error)
	Update(t *models.Theme) error
}

func (r *ThemeRepo) Get() (*models.Theme, error) {
	var t models.Theme
	if err := r.db.First(&t, 1).Error; err != nil {
		def := models.DefaultTheme
		if err := r.db.Create(&def).Error; err != nil {
			return nil, err
		}
		return &def, nil
	}
	return &t, nil
}

func (r *ThemeRepo) Update(t *models.Theme) error {
	return r.db.Model(&models.Theme{}).Where("id = 1").Updates(map[string]any{
		"light_primary": t.LightPrimary,
		"light_bg":      t.LightBg,
		"light_card":    t.LightCard,
		"light_text":    t.LightText,
		"light_border":  t.LightBorder,
		"light_danger":  t.LightDanger,
		"light_success": t.LightSuccess,
		"dark_primary":  t.DarkPrimary,
		"dark_bg":       t.DarkBg,
		"dark_card":     t.DarkCard,
		"dark_text":     t.DarkText,
		"dark_border":   t.DarkBorder,
		"dark_danger":   t.DarkDanger,
		"dark_success":  t.DarkSuccess,
	}).Error
}
