package services

import (
	"time"
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"
)

type FixtureService interface {
	CreateFixture(fixture *models.Fixture) error
	GetFixtures() ([]models.Fixture, error)
	UpdateScore(id string, homeScore, awayScore int) error
}

type fixtureService struct {
	repo repository.FixtureRepository
}

func NewFixtureService(repo repository.FixtureRepository) FixtureService {
	return &fixtureService{repo}
}

func (s *fixtureService) CreateFixture(fixture *models.Fixture) error {
	s.updateStatus(fixture)
	return s.repo.Create(fixture)
}

func (s *fixtureService) GetFixtures() ([]models.Fixture, error) {
	fixtures, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}
	// Sync statuses on retrieval for demo purposes
	for i := range fixtures {
		s.updateStatus(&fixtures[i])
	}
	return fixtures, nil
}

func (s *fixtureService) UpdateScore(id string, homeScore, awayScore int) error {
	fixture, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	fixture.HomeScore = homeScore
	fixture.AwayScore = awayScore
	fixture.Status = "completed"
	return s.repo.Update(fixture)
}

func (s *fixtureService) updateStatus(fixture *models.Fixture) {
	now := time.Now()
	matchTime := fixture.MatchDate
	matchEndTime := matchTime.Add(2 * time.Hour) // Approximate match duration

	if now.After(matchTime) && now.Before(matchEndTime) {
		fixture.Status = "live"
	} else if now.After(matchEndTime) {
		fixture.Status = "completed"
	} else {
		fixture.Status = "upcoming"
	}
}
