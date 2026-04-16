package repos

import (
	"headless-cms/internal/models"

	"gorm.io/gorm"
)

type VideoRepo struct {
	db *gorm.DB
}

func NewVideoRepo(db *gorm.DB) *VideoRepo {
	return &VideoRepo{db: db}
}

type IVideoRepo interface {
	Create(v *models.Video) error
	GetByID(id uint) (*models.Video, error)
	ListAll() ([]models.Video, error)
}

func (r *VideoRepo) Create(v *models.Video) error {
	return r.db.Create(v).Error
}

func (r *VideoRepo) GetByID(id uint) (*models.Video, error) {
	var v models.Video
	if err := r.db.Preload("ThumbnailImage").First(&v, id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *VideoRepo) ListAll() ([]models.Video, error) {
	var res []models.Video
	if err := r.db.Preload("ThumbnailImage").Order("id desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

