package models

import "github.com/google/uuid"

type Ticket struct {
	BaseModel
	FixtureID         uuid.UUID `json:"fixture_id" gorm:"type:char(36);index"`
	Fixture           Fixture   `json:"fixture" gorm:"foreignKey:FixtureID"`
	Category          string    `json:"category" gorm:"size:50"` // VIP, Regular, VVIP
	Price             float64   `json:"price"`
	AvailableQuantity int       `json:"available_quantity"`
}

type TicketOrder struct {
	BaseModel
	TicketID    uuid.UUID  `json:"ticket_id" gorm:"type:char(36);index"`
	Ticket      Ticket     `json:"ticket" gorm:"foreignKey:TicketID"`
	UserID      *uuid.UUID `json:"user_id,omitempty" gorm:"type:char(36);index"`
	User        *User      `json:"user,omitempty" gorm:"foreignKey:UserID"`
	FullName    string     `json:"full_name"`
	Mobile      string     `json:"mobile" gorm:"size:20"`
	Email       string     `json:"email"`
	Category    string     `json:"category" gorm:"size:50"`
	OrderNumber string     `json:"order_number" gorm:"type:varchar(20);uniqueIndex"`
	Quantity    int        `json:"quantity"`
	TotalAmount float64    `json:"total_amount"`
	Status      string     `json:"status" gorm:"size:20"` // pending, paid, etc.
	QRCodeURL   string     `json:"qr_code_url"`
}
