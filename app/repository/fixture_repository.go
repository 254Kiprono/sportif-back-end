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
	query := `INSERT INTO fixtures (id, created_at, updated_at, home_team, away_team, match_date, venue, home_score, away_score, status) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, fixture.ID, fixture.CreatedAt, fixture.UpdatedAt, fixture.HomeTeam, fixture.AwayTeam,
		fixture.MatchDate, fixture.Venue, fixture.HomeScore, fixture.AwayScore, fixture.Status).Error
}

func (r *fixtureRepository) GetAll() ([]models.Fixture, error) {
	var fixtures []models.Fixture
	query := `SELECT * FROM fixtures WHERE deleted_at IS NULL ORDER BY match_date ASC`
	err := r.db.Raw(query).Scan(&fixtures).Error
	return fixtures, err
}

func (r *fixtureRepository) GetByID(id string) (*models.Fixture, error) {
	var fixture models.Fixture
	query := `SELECT * FROM fixtures WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, id).Scan(&fixture).Error
	return &fixture, err
}

func (r *fixtureRepository) Update(fixture *models.Fixture) error {
	query := `UPDATE fixtures SET updated_at = ?, home_score = ?, away_score = ?, status = ?, venue = ?, match_date = ? WHERE id = ?`
	return r.db.Exec(query, fixture.UpdatedAt, fixture.HomeScore, fixture.AwayScore, fixture.Status, fixture.Venue,
		fixture.MatchDate, fixture.ID).Error
}
