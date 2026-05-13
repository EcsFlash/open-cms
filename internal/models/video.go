package models

import "time"

type Video struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Name string `json:"name" gorm:"not null"`
	Mime string `json:"mime" gorm:"not null"`

	Bucket    string  `json:"bucket" gorm:"not null"`
	ObjectKey string  `json:"object_key" gorm:"not null"`
	Duration  float64 `json:"duration_sec"`

	UploaderID uint `json:"uploader_id" gorm:"index"`
	Uploader   User `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	ThumbnailImageID *uint  `json:"thumbnail_image_id" gorm:"index"`
	ThumbnailImage   *Image `json:"thumbnail_image,omitempty" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`

	Index uint `json:"index" gorm:"column:index_priority;autoIncrement;uniqueIndex;not null;"`
}
