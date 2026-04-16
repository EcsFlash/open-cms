package services

import (
	"headless-cms/internal/models"
	"headless-cms/internal/repos"
)

type ArticleService struct {
	repo repos.IArticleRepo
}

type IArticleService interface {
	Create(a *models.Article) error
	Update(a *models.Article) error
	Delete(id uint) error
	GetByID(id uint) (*models.Article, error)
	ListAll() ([]models.Article, error)
	ListBySection(sectionID uint) ([]models.Article, error)
}

func NewArticleService(repo repos.IArticleRepo) *ArticleService {
	return &ArticleService{repo: repo}
}

func (s *ArticleService) Create(a *models.Article) error {
	return s.repo.Create(a)
}

func (s *ArticleService) Update(a *models.Article) error {
	return s.repo.Update(a)
}

func (s *ArticleService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *ArticleService) GetByID(id uint) (*models.Article, error) {
	return s.repo.GetByID(id)
}

func (s *ArticleService) ListAll() ([]models.Article, error) {
	return s.repo.ListAll()
}

func (s *ArticleService) ListBySection(sectionID uint) ([]models.Article, error) {
	return s.repo.ListBySection(sectionID)
}

