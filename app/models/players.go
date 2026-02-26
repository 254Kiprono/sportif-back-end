package models

import (
	"time"

	"gorm.io/gorm"
)

type Player struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	Name         string         `json:"name" gorm:"not null"`
	Position     string         `json:"position"`
	JerseyNumber int            `json:"jersey_number"`
	Nationality  string         `json:"nationality"`
	Age          int            `json:"age"`
	Appearances  int            `json:"appearances" gorm:"default:0"`
	Goals        int            `json:"goals" gorm:"default:0"`
	Assists      int            `json:"assists" gorm:"default:0"`
	ImageURL     string         `json:"image_url"`
}
