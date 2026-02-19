package models

import "github.com/google/uuid"

type Donation struct {
	BaseModel
	UserID        *uuid.UUID `json:"user_id" gorm:"index"` // Nullable for guests
	User          *User      `json:"user" gorm:"foreignKey:UserID"`
	Amount        float64    `json:"amount" gorm:"not null"`
	Message       string     `json:"message"`
	PaymentStatus string     `json:"payment_status" gorm:"default:'pending'"`
}
