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
	theme    services.IThemeService
	settings services.ISettingsService
}

type IUseCase interface {
	// sections
	CreateSection(actor Actor, section *models.Section) error
	UpdateSection(actor Actor, section *models.Section) error
	DeleteSection(actor Actor, id uint) error
	GetSectionByID(id uint) (*models.Section, error)
	ListSections() ([]models.Section, error)
	ListSectionChildren(parentID uint) ([]models.Section, error)

	// articles
	CreateArticle(actor Actor, a *models.Article) error
	UpdateArticle(actor Actor, a *models.Article) error
	DeleteArticle(actor Actor, id uint) error
	GetArticleByID(id uint) (*models.Article, error)
	ListArticles() ([]models.Article, error)
	ListArticlesBySection(sectionID uint) ([]models.Article, error)

	// news
	CreateNews(actor Actor, n *models.News) error
	UpdateNews(actor Actor, n *models.News) error
	DeleteNews(actor Actor, id uint) error
	GetNewsByID(id uint) (*models.News, error)
	ListNews() ([]models.News, error)

	// auth
	Register(nickname, password string) (*models.User, error)
	RegisterSupervisor(nickname, password string, role models.Role) (*models.User, error)
	Login(nickname, password string) (string, *models.User, error)
	RemoveSupervisor(nickname string) error

	// media
	UploadImage(actor Actor, ctx context.Context, file multipart.File, filename, contentType, name string) (*models.Image, error)
	UploadVideo(actor Actor, ctx context.Context, file multipart.File, filename, contentType, name string, thumbnailImageID *uint) (*models.Video, error)
	GetImageByID(id uint) (*models.Image, error)
	GetVideoByID(id uint) (*models.Video, error)
	ListImages() ([]models.Image, error)
	ListVideos() ([]models.Video, error)
	RenameImage(actor Actor, id uint, name string) error
	DeleteImage(actor Actor, ctx context.Context, id uint) error
	RenameVideo(actor Actor, id uint, name string) error
	DeleteVideo(actor Actor, ctx context.Context, id uint) error

	// theme
	GetTheme() (*models.Theme, error)
	UpdateTheme(t *models.Theme) error

	// settings
	GetSettings() (*models.Settings, error)
	UpdateSettings(s *models.Settings) error
}

func NewUseCase(
	sections services.ISectionService,
	articles services.IArticleService,
	news services.INewsService,
	auth services.IAuthService,
	media services.IMediaService,
	theme services.IThemeService,
	settings services.ISettingsService,
) *UseCase {
	return &UseCase{
		sections: sections,
		articles: articles,
		news:     news,
		auth:     auth,
		media:    media,
		theme:    theme,
		settings: settings,
	}
}

func (u *UseCase) CreateSection(actor Actor, section *models.Section) error {
	section.ID = 0
	section.AuthorID = actor.UserID
	return u.sections.Create(section)
}
func (u *UseCase) UpdateSection(actor Actor, section *models.Section) error {
	existing, err := u.sections.GetByID(section.ID)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, existing.AuthorID); err != nil {
		return err
	}
	return u.sections.Update(section)
}
func (u *UseCase) DeleteSection(actor Actor, id uint) error {
	existing, err := u.sections.GetByID(id)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, existing.AuthorID); err != nil {
		return err
	}
	return u.sections.Delete(id)
}
func (u *UseCase) GetSectionByID(id uint) (*models.Section, error) { return u.sections.GetByID(id) }
func (u *UseCase) ListSections() ([]models.Section, error)         { return u.sections.ListAll() }
func (u *UseCase) ListSectionChildren(parentID uint) ([]models.Section, error) {
	return u.sections.ListChildren(parentID)
}

func (u *UseCase) CreateArticle(actor Actor, a *models.Article) error {
	a.ID = 0
	a.AuthorID = actor.UserID
	return u.articles.Create(a)
}
func (u *UseCase) UpdateArticle(actor Actor, a *models.Article) error {
	existing, err := u.articles.GetByID(a.ID)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, existing.AuthorID); err != nil {
		return err
	}
	return u.articles.Update(a)
}
func (u *UseCase) DeleteArticle(actor Actor, id uint) error {
	existing, err := u.articles.GetByID(id)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, existing.AuthorID); err != nil {
		return err
	}
	return u.articles.Delete(id)
}
func (u *UseCase) GetArticleByID(id uint) (*models.Article, error) { return u.articles.GetByID(id) }
func (u *UseCase) ListArticles() ([]models.Article, error)         { return u.articles.ListAll() }
func (u *UseCase) ListArticlesBySection(sectionID uint) ([]models.Article, error) {
	return u.articles.ListBySection(sectionID)
}

func (u *UseCase) CreateNews(actor Actor, n *models.News) error {
	n.ID = 0
	n.AuthorID = actor.UserID
	return u.news.Create(n)
}
func (u *UseCase) UpdateNews(actor Actor, n *models.News) error {
	existing, err := u.news.GetByID(n.ID)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, existing.AuthorID); err != nil {
		return err
	}
	return u.news.Update(n)
}
func (u *UseCase) DeleteNews(actor Actor, id uint) error {
	existing, err := u.news.GetByID(id)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, existing.AuthorID); err != nil {
		return err
	}
	return u.news.Delete(id)
}
func (u *UseCase) GetNewsByID(id uint) (*models.News, error) { return u.news.GetByID(id) }
func (u *UseCase) ListNews() ([]models.News, error)          { return u.news.ListAll() }

func (u *UseCase) Register(nickname, password string) (*models.User, error) {
	return u.auth.Register(nickname, password)
}

func (u *UseCase) RegisterSupervisor(nickname, password string, role models.Role) (*models.User, error) {
	return u.auth.RegisterSupervisor(nickname, password, role)
}
func (u *UseCase) Login(nickname, password string) (string, *models.User, error) {
	return u.auth.Login(nickname, password)
}

func (u *UseCase) RemoveSupervisor(nickname string) error {
	return u.auth.RemoveSupervisor(nickname)
}

func (u *UseCase) UploadImage(actor Actor, ctx context.Context, file multipart.File, filename, contentType, name string) (*models.Image, error) {
	return u.media.UploadImage(ctx, file, filename, contentType, name, actor.UserID)
}

func (u *UseCase) UploadVideo(actor Actor, ctx context.Context, file multipart.File, filename, contentType, name string, thumbnailImageID *uint) (*models.Video, error) {
	return u.media.UploadVideo(ctx, file, filename, contentType, name, thumbnailImageID, actor.UserID)
}

func (u *UseCase) GetImageByID(id uint) (*models.Image, error) { return u.media.GetImageByID(id) }
func (u *UseCase) GetVideoByID(id uint) (*models.Video, error) { return u.media.GetVideoByID(id) }
func (u *UseCase) ListImages() ([]models.Image, error)         { return u.media.ListImages() }
func (u *UseCase) ListVideos() ([]models.Video, error)         { return u.media.ListVideos() }
func (u *UseCase) RenameImage(actor Actor, id uint, name string) error {
	img, err := u.media.GetImageByID(id)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, img.UploaderID); err != nil {
		return err
	}
	return u.media.RenameImage(id, name)
}
func (u *UseCase) DeleteImage(actor Actor, ctx context.Context, id uint) error {
	img, err := u.media.GetImageByID(id)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, img.UploaderID); err != nil {
		return err
	}
	return u.media.DeleteImage(ctx, id)
}
func (u *UseCase) RenameVideo(actor Actor, id uint, name string) error {
	v, err := u.media.GetVideoByID(id)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, v.UploaderID); err != nil {
		return err
	}
	return u.media.RenameVideo(id, name)
}
func (u *UseCase) DeleteVideo(actor Actor, ctx context.Context, id uint) error {
	v, err := u.media.GetVideoByID(id)
	if err != nil {
		return err
	}
	if err := EnsureCanMutate(actor, v.UploaderID); err != nil {
		return err
	}
	return u.media.DeleteVideo(ctx, id)
}

func (u *UseCase) GetTheme() (*models.Theme, error)  { return u.theme.Get() }
func (u *UseCase) UpdateTheme(t *models.Theme) error { return u.theme.Update(t) }

func (u *UseCase) GetSettings() (*models.Settings, error)  { return u.settings.Get() }
func (u *UseCase) UpdateSettings(s *models.Settings) error { return u.settings.Update(s) }
