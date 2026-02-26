package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"
)

type PlayerService interface {
	CreatePlayer(player *models.Player) error
	GetPlayers() ([]models.Player, error)
	GetPlayer(id uint) (*models.Player, error)
	UpdatePlayer(player *models.Player) error
	DeletePlayer(id uint) error
}

type playerService struct {
	repo repository.PlayerRepository
}

func NewPlayerService(repo repository.PlayerRepository) PlayerService {
	return &playerService{repo}
}

func (s *playerService) CreatePlayer(player *models.Player) error {
	return s.repo.Create(player)
}

func (s *playerService) GetPlayers() ([]models.Player, error) {
	return s.repo.GetAll()
}

func (s *playerService) GetPlayer(id uint) (*models.Player, error) {
	return s.repo.GetByID(id)
}

func (s *playerService) UpdatePlayer(player *models.Player) error {
	return s.repo.Update(player)
}

func (s *playerService) DeletePlayer(id uint) error {
	return s.repo.Delete(id)
}
