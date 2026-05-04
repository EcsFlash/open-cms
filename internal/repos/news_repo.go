package repos

import (
	"headless-cms/internal/models"

	"gorm.io/gorm"
)

type NewsRepo struct {
	db *gorm.DB
}

func NewNewsRepo(db *gorm.DB) *NewsRepo {
	return &NewsRepo{db: db}
}

type INewsRepo interface {
	Create(n *models.News) error
	Update(n *models.News) error
	Delete(id uint) error
	GetByID(id uint) (*models.News, error)
	ListAll() ([]models.News, error)
}

func (r *NewsRepo) Create(n *models.News) error {
	return r.db.Create(n).Error
}

func (r *NewsRepo) Update(n *models.News) error {
	return r.db.Model(&models.News{}).Where("id = ?", n.ID).Updates(map[string]any{
		"title":            n.Title,
		"preview_image_id": n.PreviewImageID,
		"short":            n.Short,
		"content": n.Content,
	}).Error
}

func (r *NewsRepo) Delete(id uint) error {
	return r.db.Delete(&models.News{}, id).Error
}

func (r *NewsRepo) GetByID(id uint) (*models.News, error) {
	var n models.News
	if err := r.db.Preload("PreviewImage").First(&n, id).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *NewsRepo) ListAll() ([]models.News, error) {
	var res []models.News
	if err := r.db.Order("id desc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}
