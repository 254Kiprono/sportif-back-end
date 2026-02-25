package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type FanRepository interface {
	GetAll() ([]models.Fan, error)
	GetByID(id string) (*models.Fan, error)
	Create(fan *models.Fan) error
	Update(fan *models.Fan) error
	Delete(id string) error
}

type fanRepository struct {
	db *gorm.DB
}

func NewFanRepository(db *gorm.DB) FanRepository {
	return &fanRepository{db}
}

func (r *fanRepository) GetAll() ([]models.Fan, error) {
	var fans []models.Fan
	err := r.db.Raw("SELECT * FROM fans ORDER BY created_at DESC").Scan(&fans).Error
	return fans, err
}

func (r *fanRepository) GetByID(id string) (*models.Fan, error) {
	var fan models.Fan
	err := r.db.Raw("SELECT * FROM fans WHERE id = ? LIMIT 1", id).Scan(&fan).Error
	return &fan, err
}

func (r *fanRepository) Create(fan *models.Fan) error {
	return r.db.Exec("INSERT INTO fans (id, created_at, updated_at, name, email, tier, join_date, location, membership_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		fan.ID, fan.CreatedAt, fan.UpdatedAt, fan.Name, fan.Email, fan.Tier, fan.JoinDate, fan.Location, fan.MembershipID).Error
}

func (r *fanRepository) Update(fan *models.Fan) error {
	return r.db.Exec("UPDATE fans SET name = ?, email = ?, tier = ?, join_date = ?, location = ?, membership_id = ?, updated_at = NOW() WHERE id = ?",
		fan.Name, fan.Email, fan.Tier, fan.JoinDate, fan.Location, fan.MembershipID, fan.ID).Error
}

func (r *fanRepository) Delete(id string) error {
	return r.db.Exec("DELETE FROM fans WHERE id = ?", id).Error
}
