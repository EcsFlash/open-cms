package models

import "time"

type Image struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name   string `json:"name" gorm:"not null"` // used as alt
	Mime   string `json:"mime" gorm:"not null"`
	Width  int    `json:"width"`
	Height int    `json:"height"`

	Bucket string `json:"bucket" gorm:"not null"`

	ObjectKeyOriginal  string `json:"object_key_original" gorm:"not null"`
	ObjectKeyLarge     string `json:"object_key_large"`
	ObjectKeyMedium    string `json:"object_key_medium"`
	ObjectKeyMini      string `json:"object_key_mini"`
	ObjectKeyThumbnail string `json:"object_key_thumbnail"`
}

