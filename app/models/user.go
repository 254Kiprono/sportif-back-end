package models

import (
	"github.com/google/uuid"
)

type User struct {
	BaseModel
	FullName string    `json:"full_name" gorm:"not null"`
	Username string    `json:"username" gorm:"uniqueIndex;not null"`
	Email    string    `json:"email" gorm:"uniqueIndex;not null"`
	Phone    string    `json:"phone" gorm:"uniqueIndex;not null"`
	Password string    `json:"-" gorm:"not null"`
	RoleID   uuid.UUID `json:"role_id" gorm:"type:uuid;not null"`
	Role     Role      `json:"role" gorm:"foreignKey:RoleID"`
}
