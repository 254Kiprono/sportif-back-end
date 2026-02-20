package models

import "github.com/google/uuid"

type News struct {
	BaseModel
	Title     string    `json:"title" gorm:"not null"`
	Content   string    `json:"content" gorm:"type:text"`
	ImageURL  string    `json:"image_url"`
	AuthorID  uuid.UUID `json:"author_id"`
	Author    User      `json:"author" gorm:"foreignKey:AuthorID"`
	Published bool      `json:"published" gorm:"default:false"`
	Status    string    `json:"status" gorm:"default:'draft'"` // draft / published
}
