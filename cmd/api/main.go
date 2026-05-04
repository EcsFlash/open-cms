package main

import (
	_ "headless-cms/docs"
	"context"
	"errors"
	"headless-cms/internal/config"
	"headless-cms/internal/db"
	"headless-cms/internal/handler"
	"headless-cms/internal/models"
	"headless-cms/internal/repos"
	"headless-cms/internal/services"
	"headless-cms/internal/usecase"
	"headless-cms/pkg/sl"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Headless CMS API
// @version 1.0
// @description Headless CMS без фронтенда. Разделы, статьи (Markdown), новости (Markdown + автогенерация краткого описания), пользователи/роли, медиа (MinIO).
// @BasePath /
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите JWT токен в формате: Bearer <token>
func main() {
	cfg := config.MustLoad()
	log := sl.SetupLogger(cfg.Env)

	gdb, err := db.Open(cfg)
	if err != nil {
		log.Error("db open failed", sl.Err(err))
		return
	}
	if err := gdb.AutoMigrate(
		&models.Section{},
		&models.Article{},
		&models.Image{},
		&models.Video{},
		&models.News{},
		&models.User{},
		&models.Theme{},
		&models.Settings{},
	); err != nil {
		log.Error("db migrate failed", sl.Err(err))
		return
	}

	sectionRepo := repos.NewSectionRepo(gdb)
	articleRepo := repos.NewArticleRepo(gdb)
	newsRepo := repos.NewNewsRepo(gdb)
	userRepo := repos.NewUserRepo(gdb)
	imageRepo := repos.NewImageRepo(gdb)
	videoRepo := repos.NewVideoRepo(gdb)

	sectionSvc := services.NewSectionService(sectionRepo)
	articleSvc := services.NewArticleService(articleRepo)
	newsSvc := services.NewNewsService(newsRepo)
	authSvc := services.NewAuthService(cfg, userRepo)

	minioCli, err := services.NewMinIO(cfg)
	if err != nil {
		log.Error("minio init failed", sl.Err(err))
		return
	}
	ctx := context.Background()
	if err := services.EnsureBucket(ctx, minioCli, cfg.MinIO.ImagesBucket); err != nil {
		log.Error("minio ensure images bucket failed", sl.Err(err))
		return
	}
	if err := services.EnsureBucket(ctx, minioCli, cfg.MinIO.VideosBucket); err != nil {
		log.Error("minio ensure videos bucket failed", sl.Err(err))
		return
	}

	mediaSvc := services.NewMediaService(cfg, minioCli, imageRepo, videoRepo)

	themeRepo := repos.NewThemeRepo(gdb)
	themeSvc := services.NewThemeService(themeRepo)

	settingsRepo := repos.NewSettingsRepo(gdb)
	settingsSvc := services.NewSettingsService(settingsRepo)

	uc := usecase.NewUseCase(sectionSvc, articleSvc, newsSvc, authSvc, mediaSvc, themeSvc, settingsSvc)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	e.Use(middleware.RequestID())
	e.Use(handler.RequestLogger(*log))
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	handler.RegisterRoutes(e, cfg, uc, *log)

	log.Info("server starting", "addr", cfg.HTTPServer.Address)
	if err := e.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server failed", sl.Err(err))
	}
}

