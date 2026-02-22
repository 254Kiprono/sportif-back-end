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
	err := r.db.First(&jersey, "id = ?", id).Error
	return &jersey, err
}

func (r *storeRepository) UpdateJersey(jersey *models.Jersey) error {
	return r.db.Save(jersey).Error
}

func (r *storeRepository) GetOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.
		Preload("User").
		Preload("Items.Product").
		Order("created_at DESC").
		Find(&orders).Error
	return orders, err
}

func (r *storeRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *storeRepository) UpdateJerseyStock(id string, quantity int) error {
	return r.db.Model(&models.Jersey{}).Where("id = ?", id).Update("stock_quantity", gorm.Expr("stock_quantity - ?", quantity)).Error
}

func (r *storeRepository) Transaction(fn func(repo StoreRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewStoreRepository(tx))
	})
}
