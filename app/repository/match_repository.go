package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type MatchRepository interface {
	// Lineup
	GetLineup(fixtureID string) (*models.Lineup, error)
	SaveLineup(lineup *models.Lineup) error
	DeleteLineup(fixtureID string) error

	// Events
	GetEvents(fixtureID string) ([]models.MatchEvent, error)
	CreateEvent(event *models.MatchEvent) error
	DeleteEvent(id string) error

	// Transactions
	Transaction(fn func(repo MatchRepository) error) error
	GetDB() *gorm.DB
}

type matchRepository struct {
	db *gorm.DB
}

func NewMatchRepository(db *gorm.DB) MatchRepository {
	return &matchRepository{db}
}

func (r *matchRepository) GetDB() *gorm.DB {
	return r.db
}

func (r *matchRepository) GetLineup(fixtureID string) (*models.Lineup, error) {
	var lineup models.Lineup
	err := r.db.Preload("Players.Player").Where("fixture_id = ?", fixtureID).First(&lineup).Error
	if err != nil {
		return nil, err
	}
	return &lineup, nil
}

func (r *matchRepository) SaveLineup(lineup *models.Lineup) error {
	lineup.Initialize()
	// Using a transaction internally to handle the clear-and-replace for players
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete existing lineup players if lineup exists
		var existing models.Lineup
		if err := tx.Where("fixture_id = ?", lineup.FixtureID).First(&existing).Error; err == nil {
			tx.Where("lineup_id = ?", existing.ID).Delete(&models.LineupPlayer{})
			lineup.ID = existing.ID // Keep same ID
			// Omit "Players" to prevent GORM from trying to insert them during Save
			if err := tx.Omit("Players").Save(lineup).Error; err != nil {
				return err
			}
		} else {
			// Omit "Players" to prevent GORM from trying to insert them during Create
			if err := tx.Omit("Players").Create(lineup).Error; err != nil {
				return err
			}
		}

		for i := range lineup.Players {
			lineup.Players[i].LineupID = lineup.ID
			lineup.Players[i].Initialize()
			// Now manually create each player record
			if err := tx.Create(&lineup.Players[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *matchRepository) DeleteLineup(fixtureID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var lineup models.Lineup
		if err := tx.Where("fixture_id = ?", fixtureID).First(&lineup).Error; err != nil {
			return err
		}
		tx.Where("lineup_id = ?", lineup.ID).Delete(&models.LineupPlayer{})
		return tx.Delete(&lineup).Error
	})
}

func (r *matchRepository) GetEvents(fixtureID string) ([]models.MatchEvent, error) {
	var events []models.MatchEvent
	err := r.db.Preload("Player").Preload("AssistPlayer").Preload("PlayerOut").
		Where("fixture_id = ?", fixtureID).Order("minute asc, created_at asc").Find(&events).Error
	return events, err
}

func (r *matchRepository) CreateEvent(event *models.MatchEvent) error {
	event.Initialize()
	return r.db.Create(event).Error
}

func (r *matchRepository) DeleteEvent(id string) error {
	return r.db.Delete(&models.MatchEvent{}, "id = ?", id).Error
}

func (r *matchRepository) Transaction(fn func(repo MatchRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewMatchRepository(tx))
	})
}
