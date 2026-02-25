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

type TicketOrder struct {
	BaseModel
	TicketID    uuid.UUID `json:"ticket_id" gorm:"index"`
	Ticket      Ticket    `json:"ticket" gorm:"foreignKey:TicketID"`
	FullName    string    `json:"full_name"`
	Mobile      string    `json:"mobile"`
	Email       string    `json:"email"`
	Category    string    `json:"category"`
	OrderNumber string    `json:"order_number" gorm:"uniqueIndex"`
	Quantity    int       `json:"quantity"`
	TotalAmount float64   `json:"total_amount"`
	Status      string    `json:"status"` // pending, paid, etc.
	QRCodeURL   string    `json:"qr_code_url"`
}
