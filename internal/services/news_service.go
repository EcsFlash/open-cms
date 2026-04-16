package services

import (
	"headless-cms/internal/models"
	"headless-cms/internal/repos"
	"regexp"
	"strings"
	"unicode/utf8"
)

type NewsService struct {
	repo repos.INewsRepo
}

type INewsService interface {
	Create(n *models.News) error
	Update(n *models.News) error
	Delete(id uint) error
	GetByID(id uint) (*models.News, error)
	ListAll() ([]models.News, error)
}

func NewNewsService(repo repos.INewsRepo) *NewsService {
	return &NewsService{repo: repo}
}

const newsShortMaxRunes = 84

func compileNewsShort(markdown string) string {
	// Remove markdown images to avoid polluted "short" field.
	imagePattern := `!\[.*?\]\(.*?\)`
	imageRe := regexp.MustCompile(imagePattern)
	withoutImages := imageRe.ReplaceAllString(markdown, "")

	// Collapse whitespace/newlines into single spaces.
	wsRe := regexp.MustCompile(`\s+`)
	normalized := strings.TrimSpace(wsRe.ReplaceAllString(withoutImages, " "))

	r := []rune(normalized)
	if utf8.RuneCountInString(normalized) > newsShortMaxRunes {
		return string(r[:newsShortMaxRunes]) + "..."
	}
	return normalized
}

func (s *NewsService) Create(n *models.News) error {
	n.Short = compileNewsShort(n.Content)
	return s.repo.Create(n)
}

func (s *NewsService) Update(n *models.News) error {
	n.Short = compileNewsShort(n.Content)
	return s.repo.Update(n)
}

func (s *NewsService) Delete(id uint) error {
	return s.repo.Delete(id)
}

func (s *NewsService) GetByID(id uint) (*models.News, error) {
	return s.repo.GetByID(id)
}

func (s *NewsService) ListAll() ([]models.News, error) {
	return s.repo.ListAll()
}
