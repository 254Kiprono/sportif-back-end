package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	GetAll() ([]models.Payment, error)
	GetByID(id string) (*models.Payment, error)
	Create(payment *models.Payment) error
	Update(payment *models.Payment) error
	Delete(id string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) GetAll() ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.Raw("SELECT * FROM payments ORDER BY created_at DESC").Scan(&payments).Error
	return payments, err
}

func (r *paymentRepository) GetByID(id string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.Raw("SELECT * FROM payments WHERE id = ? LIMIT 1", id).Scan(&payment).Error
	return &payment, err
}

func (r *paymentRepository) Create(payment *models.Payment) error {
	return r.db.Exec("INSERT INTO payments (id, created_at, updated_at, order_id, customer, amount, method, date, status) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		payment.ID, payment.CreatedAt, payment.UpdatedAt, payment.OrderID, payment.Customer, payment.Amount, payment.Method, payment.Date, payment.Status).Error
}

func (r *paymentRepository) Update(payment *models.Payment) error {
	return r.db.Exec("UPDATE payments SET order_id = ?, customer = ?, amount = ?, method = ?, date = ?, status = ?, updated_at = NOW() WHERE id = ?",
		payment.OrderID, payment.Customer, payment.Amount, payment.Method, payment.Date, payment.Status, payment.ID).Error
}

func (r *paymentRepository) Delete(id string) error {
	return r.db.Exec("DELETE FROM payments WHERE id = ?", id).Error
}
