package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type SponsorRepository interface {
	GetAll() ([]models.Sponsor, error)
	GetByID(id string) (*models.Sponsor, error)
	Create(sponsor *models.Sponsor) error
	Update(sponsor *models.Sponsor) error
	Delete(id string) error
}

type sponsorRepository struct {
	db *gorm.DB
}

func NewSponsorRepository(db *gorm.DB) SponsorRepository {
	return &sponsorRepository{db}
}

func (r *sponsorRepository) GetAll() ([]models.Sponsor, error) {
	var sponsors []models.Sponsor
	err := r.db.Order("created_at DESC").Find(&sponsors).Error
	return sponsors, err
}

func (r *sponsorRepository) GetByID(id string) (*models.Sponsor, error) {
	var sponsor models.Sponsor
	err := r.db.Where("id = ?", id).First(&sponsor).Error
	return &sponsor, err
}

func (r *sponsorRepository) Create(sponsor *models.Sponsor) error {
	return r.db.Create(sponsor).Error
}

func (r *sponsorRepository) Update(sponsor *models.Sponsor) error {
	return r.db.Model(sponsor).Select("*").Omit("CreatedAt").Updates(sponsor).Error
}

func (r *sponsorRepository) Delete(id string) error {
	return r.db.Delete(&models.Sponsor{}, "id = ?", id).Error
}
