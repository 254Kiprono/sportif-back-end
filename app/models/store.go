package models

import "github.com/google/uuid"

type Jersey struct {
	BaseModel
	Name          string  `json:"name" gorm:"not null"`
	Description   string  `json:"description"`
	Size          string  `json:"size"`
	Price         float64 `json:"price" gorm:"not null"`
	StockQuantity int     `json:"stock_quantity" gorm:"default:0"`
	ImageURL      string  `json:"image_url"`
}

type Order struct {
	BaseModel
	UserID        uuid.UUID   `json:"user_id" gorm:"index"`
	User          User        `json:"user" gorm:"foreignKey:UserID"`
	TotalAmount   float64     `json:"total_amount"`
	Status        string      `json:"status" gorm:"default:'pending'"` // pending, paid, cancelled
	PaymentMethod string      `json:"payment_method"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	BaseModel
	OrderID   uuid.UUID `json:"order_id" gorm:"index"`
	ProductID uuid.UUID `json:"product_id" gorm:"index"`
	Product   Jersey    `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}
