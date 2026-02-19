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
	query := `INSERT INTO players (id, created_at, updated_at, name, position, jersey_number, nationality, age, appearances, goals, assists, image_url) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, player.ID, player.CreatedAt, player.UpdatedAt, player.Name, player.Position, player.JerseyNumber,
		player.Nationality, player.Age, player.Appearances, player.Goals, player.Assists, player.ImageURL).Error
}

func (r *playerRepository) GetAll() ([]models.Player, error) {
	var players []models.Player
	query := `SELECT * FROM players WHERE deleted_at IS NULL`
	err := r.db.Raw(query).Scan(&players).Error
	return players, err
}

func (r *playerRepository) GetByID(id string) (*models.Player, error) {
	var player models.Player
	query := `SELECT * FROM players WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, id).Scan(&player).Error
	return &player, err
}

func (r *playerRepository) Update(player *models.Player) error {
	query := `UPDATE players SET updated_at = ?, name = ?, position = ?, jersey_number = ?, nationality = ?, age = ?, 
	          appearances = ?, goals = ?, assists = ?, image_url = ? WHERE id = ?`
	return r.db.Exec(query, player.UpdatedAt, player.Name, player.Position, player.JerseyNumber, player.Nationality,
		player.Age, player.Appearances, player.Goals, player.Assists, player.ImageURL, player.ID).Error
}

func (r *playerRepository) Delete(id string) error {
	query := `UPDATE players SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`
	return r.db.Exec(query, id).Error
}
