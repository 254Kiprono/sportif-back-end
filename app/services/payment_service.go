package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"
)

type PaymentService interface {
	GetAll() ([]models.Payment, error)
	GetByID(id string) (*models.Payment, error)
	Create(payment *models.Payment) error
	Update(payment *models.Payment) error
	Delete(id string) error
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (s *paymentService) GetAll() ([]models.Payment, error) {
	return s.repo.GetAll()
}

func (s *paymentService) GetByID(id string) (*models.Payment, error) {
	return s.repo.GetByID(id)
}

func (s *paymentService) Create(payment *models.Payment) error {
	return s.repo.Create(payment)
}

func (s *paymentService) Update(payment *models.Payment) error {
	return s.repo.Update(payment)
}

func (s *paymentService) Delete(id string) error {
	return s.repo.Delete(id)
}
