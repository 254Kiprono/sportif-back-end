package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type StoreRepository interface {
	GetJerseys() ([]models.Jersey, error)
	GetJerseyByID(id string) (*models.Jersey, error)
	UpdateJersey(jersey *models.Jersey) error
	GetOrders() ([]models.Order, error)
	CreateOrder(order *models.Order) error
	UpdateJerseyStock(id string, quantity int) error
	Transaction(fn func(repo StoreRepository) error) error
}

type storeRepository struct {
	db *gorm.DB
}

func NewStoreRepository(db *gorm.DB) StoreRepository {
	return &storeRepository{db}
}

func (r *storeRepository) GetJerseys() ([]models.Jersey, error) {
	var jerseys []models.Jersey
	err := r.db.Where("deleted_at IS NULL").Find(&jerseys).Error
	return jerseys, err
}

func (r *storeRepository) GetJerseyByID(id string) (*models.Jersey, error) {
	var jersey models.Jersey
	query := `SELECT * FROM jerseys WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, id).Scan(&jersey).Error
	return &jersey, err
}

func (r *storeRepository) UpdateJersey(jersey *models.Jersey) error {
	return r.db.Save(jersey).Error
}

func (r *storeRepository) GetOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("User").
		Preload("Items").
		Preload("Items.Product").
		Where("orders.deleted_at IS NULL").
		Order("orders.created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *storeRepository) CreateOrder(order *models.Order) error {
	// Create order
	queryOrder := `INSERT INTO orders (id, created_at, updated_at, user_id, total_amount, status, payment_method) VALUES (?, ?, ?, ?, ?, ?, ?)`
	if err := r.db.Exec(queryOrder, order.ID, order.CreatedAt, order.UpdatedAt, order.UserID, order.TotalAmount, order.Status, order.PaymentMethod).Error; err != nil {
		return err
	}

	// Create order items
	queryItem := `INSERT INTO order_items (id, created_at, updated_at, order_id, product_id, quantity, price) VALUES (?, ?, ?, ?, ?, ?, ?)`
	for i := range order.Items {
		item := &order.Items[i]
		if err := r.db.Exec(queryItem, item.ID, item.CreatedAt, item.UpdatedAt, order.ID, item.ProductID, item.Quantity, item.Price).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *storeRepository) UpdateJerseyStock(id string, quantity int) error {
	query := `UPDATE jerseys SET stock_quantity = stock_quantity - ? WHERE id = ?`
	return r.db.Exec(query, quantity, id).Error
}

func (r *storeRepository) Transaction(fn func(repo StoreRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewStoreRepository(tx))
	})
}
