package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type StoreRepository interface {
	GetJerseys() ([]models.Jersey, error)
	GetJerseyByID(id string) (*models.Jersey, error)
	CreateJersey(jersey *models.Jersey) error
	UpdateJersey(jersey *models.Jersey) error
	DeleteJersey(id string) error
	GetOrders() ([]models.Order, error)
	CreateOrder(order *models.Order) error
	UpdateOrderStatus(id string, status string) error
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
	err := r.db.Raw("SELECT * FROM jerseys WHERE deleted_at IS NULL").Scan(&jerseys).Error
	return jerseys, err
}

func (r *storeRepository) GetJerseyByID(id string) (*models.Jersey, error) {
	var jersey models.Jersey
	err := r.db.Raw("SELECT * FROM jerseys WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&jersey).Error
	return &jersey, err
}

func (r *storeRepository) CreateJersey(jersey *models.Jersey) error {
	jersey.Initialize()
	return r.db.Exec("INSERT INTO jerseys (id, created_at, updated_at, name, description, category, price, discount_percentage, stock_quantity, image_url, is_active, variants) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		jersey.ID, jersey.CreatedAt, jersey.UpdatedAt, jersey.Name, jersey.Description, jersey.Category, jersey.Price, jersey.DiscountPercentage, jersey.StockQuantity, jersey.ImageURL, jersey.IsActive, jersey.Variants).Error
}

func (r *storeRepository) UpdateJersey(jersey *models.Jersey) error {
	return r.db.Exec("UPDATE jerseys SET name = ?, description = ?, category = ?, price = ?, discount_percentage = ?, stock_quantity = ?, image_url = ?, is_active = ?, variants = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL",
		jersey.Name, jersey.Description, jersey.Category, jersey.Price, jersey.DiscountPercentage, jersey.StockQuantity, jersey.ImageURL, jersey.IsActive, jersey.Variants, jersey.ID).Error
}

func (r *storeRepository) DeleteJersey(id string) error {
	return r.db.Exec("UPDATE jerseys SET deleted_at = NOW() WHERE id = ?", id).Error
}

func (r *storeRepository) GetOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Raw("SELECT * FROM orders WHERE deleted_at IS NULL ORDER BY created_at DESC").Scan(&orders).Error
	if err == nil {
		for i := range orders {
			// Populate User
			r.db.Raw("SELECT * FROM users WHERE id = ? AND deleted_at IS NULL", orders[i].UserID).Scan(&orders[i].User)
			// Populate OrderItems
			r.db.Raw("SELECT * FROM order_items WHERE order_id = ? AND deleted_at IS NULL", orders[i].ID).Scan(&orders[i].Items)
			for j := range orders[i].Items {
				r.db.Raw("SELECT * FROM jerseys WHERE id = ? AND deleted_at IS NULL", orders[i].Items[j].ProductID).Scan(&orders[i].Items[j].Product)
			}
		}
	}
	return orders, err
}

func (r *storeRepository) CreateOrder(order *models.Order) error {
	order.Initialize()
	err := r.db.Exec("INSERT INTO orders (id, created_at, updated_at, user_id, total_amount, status, payment_method) VALUES (?, ?, ?, ?, ?, ?, ?)",
		order.ID, order.CreatedAt, order.UpdatedAt, order.UserID, order.TotalAmount, order.Status, order.PaymentMethod).Error
	if err != nil {
		return err
	}
	for i := range order.Items {
		order.Items[i].Initialize()
		order.Items[i].OrderID = order.ID
		r.db.Exec("INSERT INTO order_items (id, created_at, updated_at, order_id, product_id, quantity, price) VALUES (?, ?, ?, ?, ?, ?, ?)",
			order.Items[i].ID, order.Items[i].CreatedAt, order.Items[i].UpdatedAt, order.Items[i].OrderID, order.Items[i].ProductID, order.Items[i].Quantity, order.Items[i].Price)
	}
	return nil
}

func (r *storeRepository) UpdateOrderStatus(id string, status string) error {
	return r.db.Exec("UPDATE orders SET status = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL", status, id).Error
}

func (r *storeRepository) UpdateJerseyStock(id string, quantity int) error {
	return r.db.Exec("UPDATE jerseys SET stock_quantity = stock_quantity - ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL", quantity, id).Error
}

func (r *storeRepository) Transaction(fn func(repo StoreRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewStoreRepository(tx))
	})
}
