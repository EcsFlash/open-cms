package models

import "time"

type Section struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name     string `json:"name" gorm:"not null"`
	ParentID *uint  `json:"parent_id" gorm:"index"`

	Parent   *Section `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Children []Section `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

