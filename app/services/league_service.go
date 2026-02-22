package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"
)

type LeagueService interface {
	GetTable() ([]models.LeagueTable, error)
	UpdateEntry(id string, entry *models.LeagueTable) error
	CreateEntry(entry *models.LeagueTable) error
	DeleteEntry(id string) error
}

type leagueService struct {
	repo repository.LeagueRepository
}

func NewLeagueService(repo repository.LeagueRepository) LeagueService {
	return &leagueService{repo}
}

func (s *leagueService) GetTable() ([]models.LeagueTable, error) {
	return s.repo.GetAll()
}

func (s *leagueService) UpdateEntry(id string, entry *models.LeagueTable) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	entry.ID = existing.ID
	entry.Calculate() // Apply business rules for points and GD

	return s.repo.Update(entry)
}

func (s *leagueService) CreateEntry(entry *models.LeagueTable) error {
	entry.Calculate()
	return s.repo.Create(entry)
}

func (s *leagueService) DeleteEntry(id string) error {
	return s.repo.Delete(id)
}
