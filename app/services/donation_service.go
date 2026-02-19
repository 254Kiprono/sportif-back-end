package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"

	"github.com/google/uuid"
)

type DonationService interface {
	Donate(amount float64, message string, userID *string) (*models.Donation, error)
	GetDonations() ([]models.Donation, error)
}

type donationService struct {
	repo repository.DonationRepository
}

func NewDonationService(repo repository.DonationRepository) DonationService {
	return &donationService{repo}
}

func (s *donationService) Donate(amount float64, message string, userID *string) (*models.Donation, error) {
	donation := &models.Donation{
		Amount:  amount,
		Message: message,
	}

	if userID != nil {
		uID, _ := uuid.Parse(*userID)
		donation.UserID = &uID
	}

	err := s.repo.Create(donation)
	return donation, err
}

func (s *donationService) GetDonations() ([]models.Donation, error) {
	return s.repo.GetAll()
}
