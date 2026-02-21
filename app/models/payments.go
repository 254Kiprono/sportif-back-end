package models

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	BaseModel
	OrderID  uuid.UUID `json:"order_id" gorm:"type:char(36);index"`
	Customer string    `json:"customer"`
	Amount   float64   `json:"amount"`
	Method   string    `json:"method"`
	Date     time.Time `json:"date"`
	Status   string    `json:"status"`
}
