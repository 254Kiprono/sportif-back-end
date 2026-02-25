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
	player.Initialize()
	return r.db.Exec("INSERT INTO players (id, created_at, updated_at, name, position, jersey_number, nationality, age, appearances, goals, assists, image_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		player.ID, player.CreatedAt, player.UpdatedAt, player.Name, player.Position, player.JerseyNumber, player.Nationality, player.Age, player.Appearances, player.Goals, player.Assists, player.ImageURL).Error
}

func (r *playerRepository) GetAll() ([]models.Player, error) {
	var players []models.Player
	err := r.db.Raw("SELECT * FROM players WHERE deleted_at IS NULL").Scan(&players).Error
	return players, err
}

func (r *playerRepository) GetByID(id string) (*models.Player, error) {
	var player models.Player
	err := r.db.Raw("SELECT * FROM players WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&player).Error
	return &player, err
}

func (r *playerRepository) Update(player *models.Player) error {
	return r.db.Exec("UPDATE players SET name = ?, position = ?, jersey_number = ?, nationality = ?, age = ?, appearances = ?, goals = ?, assists = ?, image_url = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL",
		player.Name, player.Position, player.JerseyNumber, player.Nationality, player.Age, player.Appearances, player.Goals, player.Assists, player.ImageURL, player.ID).Error
}

func (r *playerRepository) Delete(id string) error {
	return r.db.Exec("UPDATE players SET deleted_at = NOW() WHERE id = ?", id).Error
}
