package models

import "time"

type Theme struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	LightPrimary string `json:"light_primary" gorm:"not null"`
	LightBg      string `json:"light_bg"      gorm:"not null"`
	LightCard    string `json:"light_card"    gorm:"not null"`
	LightText    string `json:"light_text"    gorm:"not null"`
	LightBorder  string `json:"light_border"  gorm:"not null"`
	LightDanger  string `json:"light_danger"  gorm:"not null"`
	LightSuccess string `json:"light_success" gorm:"not null"`

	DarkPrimary string `json:"dark_primary" gorm:"not null"`
	DarkBg      string `json:"dark_bg"      gorm:"not null"`
	DarkCard    string `json:"dark_card"    gorm:"not null"`
	DarkText    string `json:"dark_text"    gorm:"not null"`
	DarkBorder  string `json:"dark_border"  gorm:"not null"`
	DarkDanger  string `json:"dark_danger"  gorm:"not null"`
	DarkSuccess string `json:"dark_success" gorm:"not null"`
}

var DefaultTheme = Theme{
	ID:           1,
	LightPrimary: "#007aff",
	LightBg:      "#efeff4",
	LightCard:    "#ffffff",
	LightText:    "#1c1c1e",
	LightBorder:  "#c8c7cc",
	LightDanger:  "#ff3b30",
	LightSuccess: "#34c759",
	DarkPrimary:  "#0a84ff",
	DarkBg:       "#1c1c1e",
	DarkCard:     "#2c2c2e",
	DarkText:     "#ffffff",
	DarkBorder:   "#38383a",
	DarkDanger:   "#ff453a",
	DarkSuccess:  "#32d74b",
}
