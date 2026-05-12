package repos

import (
	"headless-cms/internal/models"

	"gorm.io/gorm"
)

type ArticleRepo struct {
	db *gorm.DB
}

func NewArticleRepo(db *gorm.DB) *ArticleRepo {
	return &ArticleRepo{db: db}
}

type IArticleRepo interface {
	Create(a *models.Article) error
	Update(a *models.Article) error
	Delete(id uint) error
	GetByID(id uint) (*models.Article, error)
	ListAll() ([]models.Article, error)
	ListBySection(sectionID uint) ([]models.Article, error)
}

func (r *ArticleRepo) Create(a *models.Article) error {
	return r.db.Create(a).Error
}

func (r *ArticleRepo) Update(a *models.Article) error {
	return r.db.Model(&models.Article{}).Where("id = ?", a.ID).Updates(map[string]any{
		"title":            a.Title,
		"content_markdown": a.ContentMarkdown,
		"section_id":       a.SectionID,
	}).Error
}

func (r *ArticleRepo) Delete(id uint) error {
	return r.db.Delete(&models.Article{}, id).Error
}

func (r *ArticleRepo) GetByID(id uint) (*models.Article, error) {
	var a models.Article
	if err := r.db.First(&a, id).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *ArticleRepo) ListAll() ([]models.Article, error) {
	var res []models.Article
	if err := r.db.Order("id asc").Find(&res).Error; err != nil {
		return nil, err
	}
	return res, nil
}

func (r *ArticleRepo) ListBySection(sectionID uint) ([]models.Article, error) {
	var res []models.Article
	if err := r.db.Order("index_priority asc").Find(&res, "section_id = ?", sectionID).Error; err != nil {
		return nil, err
	}
	return res, nil
}
