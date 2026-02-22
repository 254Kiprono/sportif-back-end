package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type LeagueRepository interface {
	GetAll() ([]models.LeagueTable, error)
	GetByID(id string) (*models.LeagueTable, error)
	Update(entry *models.LeagueTable) error
	Create(entry *models.LeagueTable) error
}

type leagueRepository struct {
	db *gorm.DB
}

func NewLeagueRepository(db *gorm.DB) LeagueRepository {
	return &leagueRepository{db}
}

func (r *leagueRepository) GetAll() ([]models.LeagueTable, error) {
	var table []models.LeagueTable
	err := r.db.Order("points DESC, goal_difference DESC").Find(&table).Error
	return table, err
}

func (r *leagueRepository) GetByID(id string) (*models.LeagueTable, error) {
	var entry models.LeagueTable
	err := r.db.First(&entry, "id = ?", id).Error
	return &entry, err
}

func (r *leagueRepository) Update(entry *models.LeagueTable) error {
	return r.db.Save(entry).Error
}

func (r *leagueRepository) Create(entry *models.LeagueTable) error {
	return r.db.Create(entry).Error
}
