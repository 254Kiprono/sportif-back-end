package models

import "github.com/google/uuid"

type MembershipPlan struct {
	BaseModel
	Name           string  `json:"name" gorm:"not null"`
	Description    string  `json:"description"`
	Price          float64 `json:"price"`
	DurationMonths int     `json:"duration_months"`
	Benefits       string  `json:"benefits"`
}

type MembershipOrder struct {
	BaseModel
	UserID uuid.UUID      `json:"user_id" gorm:"type:char(36);index"`
	User   User           `json:"user" gorm:"foreignKey:UserID"`
	PlanID uuid.UUID      `json:"plan_id" gorm:"type:char(36);index"`
	Plan   MembershipPlan `json:"plan" gorm:"foreignKey:PlanID"`
	Amount float64        `json:"amount"`
	Status string         `json:"status" gorm:"default:'pending'"`
}
