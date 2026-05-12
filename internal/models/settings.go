package models

import "time"

type Settings struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ProjectName     string  `json:"project_name"     gorm:"not null;default:'CMS'"`
	LogoURL         *string `json:"logo_url"`
	FaviconURL      *string `json:"favicon_url"`
	MetaTitle       string  `json:"meta_title"       gorm:"not null;default:''"`
	MetaDescription string  `json:"meta_description" gorm:"not null;default:''"`
	MetaKeywords    string  `json:"meta_keywords"    gorm:"not null;default:''"`
	OgImageURL      *string `json:"og_image_url"`
	BannerText      string  `json:"banner_text"      gorm:"not null;default:''"`
	BannerImageURL  *string `json:"banner_image_url"`
}

var DefaultSettings = Settings{
	ID:              1,
	ProjectName:     "CMS",
	LogoURL:         nil,
	FaviconURL:      nil,
	MetaTitle:       "",
	MetaDescription: "",
	MetaKeywords:    "",
	OgImageURL:      nil,
	BannerText:      "",
	BannerImageURL:  nil,
}
