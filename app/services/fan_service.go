package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"
)

type FanService interface {
	GetAll() ([]models.Fan, error)
	GetByID(id string) (*models.Fan, error)
	Create(fan *models.Fan) error
	Update(fan *models.Fan) error
	Delete(id string) error
}

type fanService struct {
	repo repository.FanRepository
}

func NewFanService(repo repository.FanRepository) FanService {
	return &fanService{repo}
}

func (s *fanService) GetAll() ([]models.Fan, error) {
	return s.repo.GetAll()
}

func (s *fanService) GetByID(id string) (*models.Fan, error) {
	return s.repo.GetByID(id)
}

func (s *fanService) Create(fan *models.Fan) error {
	return s.repo.Create(fan)
}

func (s *fanService) Update(fan *models.Fan) error {
	return s.repo.Update(fan)
}

func (s *fanService) Delete(id string) error {
	return s.repo.Delete(id)
}
