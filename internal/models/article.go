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
}

