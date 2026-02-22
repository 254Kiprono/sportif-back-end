package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type PlayerRepository interface {
	Create(player *models.Player) error
	GetAll() ([]models.Player, error)
	GetByID(id string) (*models.Player, error)
	Update(player *models.Player) error
	Delete(id string) error
}

type playerRepository struct {
	db *gorm.DB
}

func NewPlayerRepository(db *gorm.DB) PlayerRepository {
	return &playerRepository{db}
}

func (r *playerRepository) Create(player *models.Player) error {
	return r.db.Create(player).Error
}

func (r *playerRepository) GetAll() ([]models.Player, error) {
	var players []models.Player
	err := r.db.Find(&players).Error
	return players, err
}

func (r *playerRepository) GetByID(id string) (*models.Player, error) {
	var player models.Player
	err := r.db.First(&player, "id = ?", id).Error
	return &player, err
}

func (r *playerRepository) Update(player *models.Player) error {
	return r.db.Save(player).Error
}

func (r *playerRepository) Delete(id string) error {
	return r.db.Delete(&models.Player{}, "id = ?", id).Error
}
