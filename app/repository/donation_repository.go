package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type DonationRepository interface {
	Create(donation *models.Donation) error
	GetAll() ([]models.Donation, error)
}

type donationRepository struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) DonationRepository {
	return &donationRepository{db}
}

func (r *donationRepository) Create(donation *models.Donation) error {
	query := `INSERT INTO donations (id, created_at, updated_at, user_id, amount, message, payment_status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, donation.ID, donation.CreatedAt, donation.UpdatedAt, donation.UserID, donation.Amount,
		donation.Message, donation.PaymentStatus).Error
}

func (r *donationRepository) GetAll() ([]models.Donation, error) {
	var donations []models.Donation
	query := `SELECT * FROM donations WHERE deleted_at IS NULL ORDER BY created_at DESC`
	err := r.db.Raw(query).Scan(&donations).Error
	return donations, err
}
