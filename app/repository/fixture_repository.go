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
	fixture.Initialize()
	return r.db.Exec("INSERT INTO fixtures (id, created_at, updated_at, home_team, away_team, match_date, venue, home_score, away_score, status, preview_image, preview_caption, match_photos) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		fixture.ID, fixture.CreatedAt, fixture.UpdatedAt, fixture.HomeTeam, fixture.AwayTeam, fixture.MatchDate, fixture.Venue, fixture.HomeScore, fixture.AwayScore, fixture.Status, fixture.PreviewImage, fixture.PreviewCaption, fixture.MatchPhotos).Error
}

func (r *fixtureRepository) GetAll() ([]models.Fixture, error) {
	var fixtures []models.Fixture
	err := r.db.Raw("SELECT * FROM fixtures WHERE deleted_at IS NULL ORDER BY match_date ASC").Scan(&fixtures).Error
	return fixtures, err
}

func (r *fixtureRepository) GetByID(id string) (*models.Fixture, error) {
	var fixture models.Fixture
	err := r.db.Raw("SELECT * FROM fixtures WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&fixture).Error
	return &fixture, err
}

func (r *fixtureRepository) Update(fixture *models.Fixture) error {
	return r.db.Exec("UPDATE fixtures SET home_team = ?, away_team = ?, match_date = ?, venue = ?, home_score = ?, away_score = ?, status = ?, preview_image = ?, preview_caption = ?, match_photos = ?, updated_at = NOW() WHERE id = ? AND deleted_at IS NULL",
		fixture.HomeTeam, fixture.AwayTeam, fixture.MatchDate, fixture.Venue, fixture.HomeScore, fixture.AwayScore, fixture.Status, fixture.PreviewImage, fixture.PreviewCaption, fixture.MatchPhotos, fixture.ID).Error
}

func (r *fixtureRepository) Delete(id string) error {
	return r.db.Exec("UPDATE fixtures SET deleted_at = NOW() WHERE id = ?", id).Error
}
