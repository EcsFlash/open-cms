package repos

import (
	"headless-cms/internal/models"

	"gorm.io/gorm"
)

type ImageRepo struct {
	db *gorm.DB
}

func NewImageRepo(db *gorm.DB) *ImageRepo {
	return &ImageRepo{db: db}
}

type IImageRepo interface {
	Create(img *models.Image) error
	GetByID(id uint) (*models.Image, error)
	ListAll() ([]models.Image, error)
	UpdateName(id uint, name string) error
	Delete(id uint) error
}

func (r *ImageRepo) Create(img *models.Image) error {
	return r.db.Create(img).Error
}

func (r *ImageRepo) GetByID(id uint) (*models.Image, error) {
	var img models.Image
	if err := r.db.First(&img, id).Error; err != nil {
		return nil, err
	}
	return &img, nil
}

func (r *ImageRepo) ListAll() ([]models.Image, error) {
	var res []models.Image
	if err := r.db.Order("id desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ImageRepo) UpdateName(id uint, name string) error {
	return r.db.Model(&models.Image{}).Where("id = ?", id).Update("name", name).Error
}

func (r *ImageRepo) Delete(id uint) error {
	return r.db.Delete(&models.Image{}, id).Error
}

