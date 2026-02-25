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
	Delete(id string) error
}

type leagueRepository struct {
	db *gorm.DB
}

func NewLeagueRepository(db *gorm.DB) LeagueRepository {
	return &leagueRepository{db}
}

func (r *leagueRepository) GetAll() ([]models.LeagueTable, error) {
	var table []models.LeagueTable
	err := r.db.Raw("SELECT * FROM league_tables ORDER BY points DESC, goal_difference DESC").Scan(&table).Error
	return table, err
}

func (r *leagueRepository) GetByID(id string) (*models.LeagueTable, error) {
	var entry models.LeagueTable
	err := r.db.Raw("SELECT * FROM league_tables WHERE id = ? LIMIT 1", id).Scan(&entry).Error
	return &entry, err
}

func (r *leagueRepository) Update(entry *models.LeagueTable) error {
	return r.db.Exec("UPDATE league_tables SET team_name = ?, played = ?, wins = ?, draws = ?, losses = ?, goals_for = ?, goals_against = ?, goal_difference = ?, points = ?, position = ?, updated_at = NOW() WHERE id = ?",
		entry.TeamName, entry.Played, entry.Wins, entry.Draws, entry.Losses, entry.GoalsFor, entry.GoalsAgainst, entry.GoalDifference, entry.Points, entry.Position, entry.ID).Error
}

func (r *leagueRepository) Create(entry *models.LeagueTable) error {
	return r.db.Exec("INSERT INTO league_tables (id, created_at, updated_at, team_name, played, wins, draws, losses, goals_for, goals_against, goal_difference, points, position) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		entry.ID, entry.CreatedAt, entry.UpdatedAt, entry.TeamName, entry.Played, entry.Wins, entry.Draws, entry.Losses, entry.GoalsFor, entry.GoalsAgainst, entry.GoalDifference, entry.Points, entry.Position).Error
}

func (r *leagueRepository) Delete(id string) error {
	return r.db.Exec("DELETE FROM league_tables WHERE id = ?", id).Error
}
