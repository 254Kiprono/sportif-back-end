package models

import "github.com/google/uuid"

type Jersey struct {
	BaseModel
	Name               string  `json:"name" gorm:"not null"`
	Description        string  `json:"description"`
	Category           string  `json:"category"`
	Price              float64 `json:"price" gorm:"not null"`
	DiscountPercentage float64 `json:"discount_percentage" gorm:"default:0"`
	StockQuantity      int     `json:"stock_quantity" gorm:"default:0"`
	ImageURL           string  `json:"image_url"`
	IsActive           bool    `json:"is_active" gorm:"default:true"`
	Variants           string  `json:"variants" gorm:"type:text"` // JSON string of variants
}

type Order struct {
	BaseModel
	UserID        uuid.UUID   `json:"user_id" gorm:"type:char(36);index"`
	User          User        `json:"user" gorm:"foreignKey:UserID"`
	TotalAmount   float64     `json:"total_amount"`
	Status        string      `json:"status" gorm:"default:'pending'"` // pending, paid, cancelled
	PaymentMethod string      `json:"payment_method"`
	Items         []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	BaseModel
	OrderID   uuid.UUID `json:"order_id" gorm:"type:char(36);index"`
	ProductID uuid.UUID `json:"product_id" gorm:"type:char(36);index"`
	Product   Jersey    `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int       `json:"quantity"`
	Price     float64   `json:"price"`
}
