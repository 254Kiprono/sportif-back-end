package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type FixtureRepository interface {
	Create(fixture *models.Fixture) error
	GetAll() ([]models.Fixture, error)
	GetByID(id string) (*models.Fixture, error)
	Update(fixture *models.Fixture) error
}

type fixtureRepository struct {
	db *gorm.DB
}

func NewFixtureRepository(db *gorm.DB) FixtureRepository {
	return &fixtureRepository{db}
}

func (r *fixtureRepository) Create(fixture *models.Fixture) error {
	return r.db.Create(fixture).Error
}

func (r *fixtureRepository) GetAll() ([]models.Fixture, error) {
	var fixtures []models.Fixture
	err := r.db.Order("match_date ASC").Find(&fixtures).Error
	return fixtures, err
}

func (r *fixtureRepository) GetByID(id string) (*models.Fixture, error) {
	var fixture models.Fixture
	err := r.db.First(&fixture, "id = ?", id).Error
	return &fixture, err
}

func (r *fixtureRepository) Update(fixture *models.Fixture) error {
	return r.db.Save(fixture).Error
}
