package models

import "time"

type News struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title string `json:"title" gorm:"not null"`

	PreviewImageID *uint  `json:"preview_image_id" gorm:"index"`
	PreviewImage   *Image `json:"preview_image,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Short   string `json:"short" gorm:"not null"`
	Content string `json:"content" gorm:"type:text;not null"`

	AuthorID uint `json:"author_id" gorm:"index;not null"`
	Author   User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ShowAuthor bool `json:"show_author" gorm:"default:false;not null"`
}
