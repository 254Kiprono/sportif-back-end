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
	return r.db.Create(donation).Error
}

func (r *donationRepository) GetAll() ([]models.Donation, error) {
	var donations []models.Donation
	err := r.db.Order("created_at DESC").Find(&donations).Error
	return donations, err
}
