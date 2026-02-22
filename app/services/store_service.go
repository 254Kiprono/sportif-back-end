package services

import (
	"errors"
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"

	"github.com/google/uuid"
)

type StoreService interface {
	PlaceOrder(userID string, items []models.OrderItem) (*models.Order, error)
	GetJerseys() ([]models.Jersey, error)
	CreateJersey(jersey *models.Jersey) error
	GetOrders() ([]models.Order, error)
	UpdateOrderStatus(id string, status string) error
	UpdateJersey(jersey *models.Jersey) error
	DeleteJersey(id string) error
}

type storeService struct {
	repo repository.StoreRepository
}

func NewStoreService(repo repository.StoreRepository) StoreService {
	return &storeService{repo}
}

func (s *storeService) GetJerseys() ([]models.Jersey, error) {
	return s.repo.GetJerseys()
}

func (s *storeService) CreateJersey(jersey *models.Jersey) error {
	return s.repo.CreateJersey(jersey)
}

func (s *storeService) GetOrders() ([]models.Order, error) {
	return s.repo.GetOrders()
}

func (s *storeService) UpdateOrderStatus(id string, status string) error {
	return s.repo.UpdateOrderStatus(id, status)
}

func (s *storeService) UpdateJersey(jersey *models.Jersey) error {
	return s.repo.UpdateJersey(jersey)
}

func (s *storeService) DeleteJersey(id string) error {
	return s.repo.DeleteJersey(id)
}

func (s *storeService) PlaceOrder(userID string, items []models.OrderItem) (*models.Order, error) {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	order := &models.Order{
		UserID: uID,
		Status: "pending",
	}

	var totalAmount float64

	err = s.repo.Transaction(func(txRepo repository.StoreRepository) error {
		for i := range items {
			jersey, err := txRepo.GetJerseyByID(items[i].ProductID.String())
			if err != nil {
				return errors.New("product not found")
			}

			if jersey.StockQuantity < items[i].Quantity {
				return errors.New("insufficient stock for product: " + jersey.Name)
			}

			items[i].Price = jersey.Price
			totalAmount += jersey.Price * float64(items[i].Quantity)

			// Deduct stock
			if err := txRepo.UpdateJerseyStock(jersey.ID.String(), items[i].Quantity); err != nil {
				return err
			}
		}

		order.TotalAmount = totalAmount
		order.Items = items

		return txRepo.CreateOrder(order)
	})

	return order, err
}
