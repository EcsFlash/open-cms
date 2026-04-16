package usecase

import (
	"context"
	"headless-cms/internal/models"
	"headless-cms/internal/services"
	"mime/multipart"
)

type UseCase struct {
	sections services.ISectionService
	articles services.IArticleService
	news     services.INewsService
	auth     services.IAuthService
	media    services.IMediaService
}

type IUseCase interface {
	// sections
	CreateSection(section *models.Section) error
	UpdateSection(section *models.Section) error
	DeleteSection(id uint) error
	GetSectionByID(id uint) (*models.Section, error)
	ListSections() ([]models.Section, error)
	ListSectionChildren(parentID uint) ([]models.Section, error)

	// articles
	CreateArticle(a *models.Article) error
	UpdateArticle(a *models.Article) error
	DeleteArticle(id uint) error
	GetArticleByID(id uint) (*models.Article, error)
	ListArticles() ([]models.Article, error)
	ListArticlesBySection(sectionID uint) ([]models.Article, error)

	// news
	CreateNews(n *models.News) error
	UpdateNews(n *models.News) error
	DeleteNews(id uint) error
	GetNewsByID(id uint) (*models.News, error)
	ListNews() ([]models.News, error)

	// auth
	Register(nickname, password string) (*models.User, error)
	Login(nickname, password string) (string, *models.User, error)

	// media
	UploadImage(ctx context.Context, file multipart.File, filename, contentType, name string) (*models.Image, error)
	UploadVideo(ctx context.Context, file multipart.File, filename, contentType, name string, thumbnailImageID *uint) (*models.Video, error)
	GetImageByID(id uint) (*models.Image, error)
	GetVideoByID(id uint) (*models.Video, error)
	ListImages() ([]models.Image, error)
	ListVideos() ([]models.Video, error)
}

func NewUseCase(
	sections services.ISectionService,
	articles services.IArticleService,
	news services.INewsService,
	auth services.IAuthService,
	media services.IMediaService,
) *UseCase {
	return &UseCase{
		sections: sections,
		articles: articles,
		news:     news,
		auth:     auth,
		media:    media,
	}
}

func (u *UseCase) CreateSection(section *models.Section) error             { return u.sections.Create(section) }
func (u *UseCase) UpdateSection(section *models.Section) error             { return u.sections.Update(section) }
func (u *UseCase) DeleteSection(id uint) error                             { return u.sections.Delete(id) }
func (u *UseCase) GetSectionByID(id uint) (*models.Section, error)         { return u.sections.GetByID(id) }
func (u *UseCase) ListSections() ([]models.Section, error)                 { return u.sections.ListAll() }
func (u *UseCase) ListSectionChildren(parentID uint) ([]models.Section, error) {
	return u.sections.ListChildren(parentID)
}

func (u *UseCase) CreateArticle(a *models.Article) error                   { return u.articles.Create(a) }
func (u *UseCase) UpdateArticle(a *models.Article) error                   { return u.articles.Update(a) }
func (u *UseCase) DeleteArticle(id uint) error                             { return u.articles.Delete(id) }
func (u *UseCase) GetArticleByID(id uint) (*models.Article, error)         { return u.articles.GetByID(id) }
func (u *UseCase) ListArticles() ([]models.Article, error)                 { return u.articles.ListAll() }
func (u *UseCase) ListArticlesBySection(sectionID uint) ([]models.Article, error) {
	return u.articles.ListBySection(sectionID)
}

func (u *UseCase) CreateNews(n *models.News) error                         { return u.news.Create(n) }
func (u *UseCase) UpdateNews(n *models.News) error                         { return u.news.Update(n) }
func (u *UseCase) DeleteNews(id uint) error                                { return u.news.Delete(id) }
func (u *UseCase) GetNewsByID(id uint) (*models.News, error)               { return u.news.GetByID(id) }
func (u *UseCase) ListNews() ([]models.News, error)                        { return u.news.ListAll() }

func (u *UseCase) Register(nickname, password string) (*models.User, error) {
	return u.auth.Register(nickname, password)
}

func (u *UseCase) Login(nickname, password string) (string, *models.User, error) {
	return u.auth.Login(nickname, password)
}

func (u *UseCase) UploadImage(ctx context.Context, file multipart.File, filename, contentType, name string) (*models.Image, error) {
	return u.media.UploadImage(ctx, file, filename, contentType, name)
}

func (u *UseCase) UploadVideo(ctx context.Context, file multipart.File, filename, contentType, name string, thumbnailImageID *uint) (*models.Video, error) {
	return u.media.UploadVideo(ctx, file, filename, contentType, name, thumbnailImageID)
}

func (u *UseCase) GetImageByID(id uint) (*models.Image, error) { return u.media.GetImageByID(id) }
func (u *UseCase) GetVideoByID(id uint) (*models.Video, error) { return u.media.GetVideoByID(id) }
func (u *UseCase) ListImages() ([]models.Image, error)         { return u.media.ListImages() }
func (u *UseCase) ListVideos() ([]models.Video, error)         { return u.media.ListVideos() }

