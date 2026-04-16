package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Nickname     string `json:"nickname" gorm:"not null;uniqueIndex"`
	PasswordHash string `json:"-" gorm:"not null"`
	Role         Role   `json:"role" gorm:"type:varchar(16);not null;default:'user'"`
}

