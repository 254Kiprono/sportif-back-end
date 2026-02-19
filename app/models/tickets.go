package models

import "github.com/google/uuid"

type Ticket struct {
	BaseModel
	FixtureID         uuid.UUID `json:"fixture_id" gorm:"index"`
	Fixture           Fixture   `json:"fixture" gorm:"foreignKey:FixtureID"`
	Category          string    `json:"category"` // VIP, Regular, VVIP
	Price             float64   `json:"price"`
	AvailableQuantity int       `json:"available_quantity"`
}
