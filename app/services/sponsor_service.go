package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"
)

type SponsorService interface {
	GetAll() ([]models.Sponsor, error)
	GetByID(id string) (*models.Sponsor, error)
	Create(sponsor *models.Sponsor) error
	Update(sponsor *models.Sponsor) error
	Delete(id string) error
}

type sponsorService struct {
	repo repository.SponsorRepository
}

func NewSponsorService(repo repository.SponsorRepository) SponsorService {
	return &sponsorService{repo}
}

func (s *sponsorService) GetAll() ([]models.Sponsor, error) {
	return s.repo.GetAll()
}

func (s *sponsorService) GetByID(id string) (*models.Sponsor, error) {
	return s.repo.GetByID(id)
}

func (s *sponsorService) Create(sponsor *models.Sponsor) error {
	return s.repo.Create(sponsor)
}

func (s *sponsorService) Update(sponsor *models.Sponsor) error {
	return s.repo.Update(sponsor)
}

func (s *sponsorService) Delete(id string) error {
	return s.repo.Delete(id)
}
