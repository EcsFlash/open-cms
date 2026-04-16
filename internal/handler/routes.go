package handler

import (
	"headless-cms/internal/config"
	"headless-cms/internal/models"
	"headless-cms/internal/usecase"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, cfg *config.Config, uc usecase.IUseCase, log slog.Logger) {
	h := New(cfg, uc, log)

	e.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]any{
			"status": "ok",
			"env":    cfg.Env,
		})
	})

	api := e.Group("/api")
	v1 := api.Group("/v1")

	auth := v1.Group("/auth")
	auth.POST("/register", h.Register)
	auth.POST("/login", h.Login)

	mod := v1.Group("", JWTAuth(cfg), RequireRole(models.RoleModerator, models.RoleAdmin))
	admin := v1.Group("", JWTAuth(cfg), RequireRole(models.RoleAdmin))

	// Sections
	v1.GET("/sections", h.ListSections)
	v1.GET("/sections/:id", h.GetSection)
	v1.GET("/sections/:id/children", h.ListSectionChildren)
	admin.POST("/sections", h.CreateSection)
	mod.PATCH("/sections/:id", h.UpdateSection)
	admin.DELETE("/sections/:id", h.DeleteSection)

	// Articles
	v1.GET("/articles", h.ListArticles)
	v1.GET("/articles/:id", h.GetArticle)
	v1.GET("/sections/:id/articles", h.ListArticlesBySection)
	admin.POST("/articles", h.CreateArticle)
	mod.PATCH("/articles/:id", h.UpdateArticle)
	admin.DELETE("/articles/:id", h.DeleteArticle)

	// News
	v1.GET("/news", h.ListNews)
	v1.GET("/news/:id", h.GetNews)
	admin.POST("/news", h.CreateNews)
	mod.PATCH("/news/:id", h.UpdateNews)
	admin.DELETE("/news/:id", h.DeleteNews)

	// Media
	v1.GET("/media", h.ListMedia)
	v1.GET("/images/:id", h.GetImage)
	v1.GET("/videos/:id", h.GetVideo)
	mod.POST("/media/images", h.UploadImage)
	mod.POST("/media/videos", h.UploadVideo)
}

