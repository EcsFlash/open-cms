package models

import "time"

type Article struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title           string `json:"title" gorm:"not null"`
	ContentMarkdown string `json:"content_markdown" gorm:"type:text;not null"`

	SectionID uint    `json:"section_id" gorm:"index;not null"`
	Section   Section `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	AuthorID uint `json:"author_id" gorm:"index;not null"`
	Author   User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ShowAuthor bool `json:"show_author" gorm:"default:false;not null"`

	IsVisible bool `json:"is_visible" gorm:"default:true;not null"`
	Index     uint `json:"index" gorm:"column: index_priority; default:0;not null"`
}
