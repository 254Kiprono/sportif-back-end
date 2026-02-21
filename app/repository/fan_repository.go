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
	err := r.db.Order("created_at DESC").Find(&fans).Error
	return fans, err
}

func (r *fanRepository) GetByID(id string) (*models.Fan, error) {
	var fan models.Fan
	err := r.db.Where("id = ?", id).First(&fan).Error
	return &fan, err
}

func (r *fanRepository) Create(fan *models.Fan) error {
	return r.db.Create(fan).Error
}

func (r *fanRepository) Update(fan *models.Fan) error {
	return r.db.Save(fan).Error
}

func (r *fanRepository) Delete(id string) error {
	return r.db.Delete(&models.Fan{}, "id = ?", id).Error
}
