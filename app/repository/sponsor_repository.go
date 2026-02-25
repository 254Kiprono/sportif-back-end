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
	err := r.db.Raw("SELECT * FROM sponsors ORDER BY created_at DESC").Scan(&sponsors).Error
	return sponsors, err
}

func (r *sponsorRepository) GetByID(id string) (*models.Sponsor, error) {
	var sponsor models.Sponsor
	err := r.db.Raw("SELECT * FROM sponsors WHERE id = ? LIMIT 1", id).Scan(&sponsor).Error
	return &sponsor, err
}

func (r *sponsorRepository) Create(sponsor *models.Sponsor) error {
	return r.db.Exec("INSERT INTO sponsors (id, created_at, updated_at, name, logo, website, tier, contract_end, active) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		sponsor.ID, sponsor.CreatedAt, sponsor.UpdatedAt, sponsor.Name, sponsor.Logo, sponsor.Website, sponsor.Tier, sponsor.ContractEnd, sponsor.Active).Error
}

func (r *sponsorRepository) Update(sponsor *models.Sponsor) error {
	return r.db.Exec("UPDATE sponsors SET name = ?, logo = ?, website = ?, tier = ?, contract_end = ?, active = ?, updated_at = NOW() WHERE id = ?",
		sponsor.Name, sponsor.Logo, sponsor.Website, sponsor.Tier, sponsor.ContractEnd, sponsor.Active, sponsor.ID).Error
}

func (r *sponsorRepository) Delete(id string) error {
	return r.db.Exec("DELETE FROM sponsors WHERE id = ?", id).Error
}
