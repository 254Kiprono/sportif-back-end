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
	query := `SELECT * FROM league_tables WHERE deleted_at IS NULL ORDER BY points DESC, goal_difference DESC`
	err := r.db.Raw(query).Scan(&table).Error
	return table, err
}

func (r *leagueRepository) GetByID(id string) (*models.LeagueTable, error) {
	var entry models.LeagueTable
	query := `SELECT * FROM league_tables WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, id).Scan(&entry).Error
	return &entry, err
}

func (r *leagueRepository) Update(entry *models.LeagueTable) error {
	query := `UPDATE league_tables SET updated_at = ?, team_name = ?, played = ?, wins = ?, draws = ?, losses = ?, 
	          goals_for = ?, goals_against = ?, goal_difference = ?, points = ?, position = ? WHERE id = ?`
	return r.db.Exec(query, entry.UpdatedAt, entry.TeamName, entry.Played, entry.Wins, entry.Draws, entry.Losses,
		entry.GoalsFor, entry.GoalsAgainst, entry.GoalDifference, entry.Points, entry.Position, entry.ID).Error
}

func (r *leagueRepository) Create(entry *models.LeagueTable) error {
	query := `INSERT INTO league_tables (id, created_at, updated_at, team_name, played, wins, draws, losses, 
	          goals_for, goals_against, goal_difference, points, position) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, entry.ID, entry.CreatedAt, entry.UpdatedAt, entry.TeamName, entry.Played, entry.Wins,
		entry.Draws, entry.Losses, entry.GoalsFor, entry.GoalsAgainst, entry.GoalDifference, entry.Points, entry.Position).Error
}
