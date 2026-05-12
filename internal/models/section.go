package models

import "time"

type Section struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name     string `json:"name" gorm:"not null"`
	ParentID *uint  `json:"parent_id" gorm:"index"`

	Parent   *Section  `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Children []Section `json:"children,omitempty" gorm:"foreignKey:ParentID"`

	AuthorID uint `json:"author_id" gorm:"index;not null"`
	Author   User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;"`

	ShowAuthor bool `json:"show_author" gorm:"default:false;not null"`

	IsVisible bool `json:"is_visible" gorm:"default:true;not null"`
	Index     uint `json:"index" gorm:"column: index_priority; default:0;not null"`
}
