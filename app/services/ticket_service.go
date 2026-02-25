package services

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"

	"github.com/skip2/go-qrcode"
)

type TicketService interface {
	GetTickets() ([]models.Ticket, error)
	PurchaseTicket(ticketID string, userID string, quantity int) error
	PurchaseTicketGuest(order *models.TicketOrder) error
	CreateTicket(ticket *models.Ticket) error
	DeleteTicket(id string) error
	GetAllOrders() ([]models.TicketOrder, error)
}

func (s *ticketService) DeleteTicket(id string) error {
	return s.repo.Delete(id)
}

func (s *ticketService) GetAllOrders() ([]models.TicketOrder, error) {
	return s.repo.GetAllOrders()
}

type ticketService struct {
	repo       repository.TicketRepository
	userRepo   repository.UserRepository
	storageSvc StorageService
}

func NewTicketService(repo repository.TicketRepository, userRepo repository.UserRepository, storageSvc StorageService) TicketService {
	return &ticketService{repo, userRepo, storageSvc}
}

func (s *ticketService) GetTickets() ([]models.Ticket, error) {
	return s.repo.GetAll()
}

func (s *ticketService) PurchaseTicket(ticketID string, userID string, quantity int) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil || user.ID == (models.BaseModel{}).ID {
		return errors.New("user not found")
	}

	return s.repo.Transaction(func(txRepo repository.TicketRepository) error {
		ticket, err := txRepo.GetByID(ticketID)
		if err != nil {
			return errors.New("ticket not found")
		}

		if ticket.AvailableQuantity < quantity {
			return errors.New("not enough tickets available")
		}

		uID := user.ID
		order := &models.TicketOrder{
			TicketID:    ticket.ID,
			UserID:      &uID,
			FullName:    user.FullName,
			Email:       user.Email,
			Mobile:      user.Phone,
			Category:    ticket.Category,
			Quantity:    quantity,
			TotalAmount: ticket.Price * float64(quantity),
			Status:      "paid", // Assume paid for authenticated purchases for now
		}

		orderNumber, err := generateOrderNumber()
		if err != nil {
			return err
		}
		order.OrderNumber = orderNumber

		// Generate TICKET QR Code
		if s.storageSvc != nil {
			qrContent := fmt.Sprintf("https://webuye-sportif-fc.devsinkenya.com/verify-ticket?order=%s", order.OrderNumber)
			png, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
			if err == nil {
				fileName := fmt.Sprintf("ticket_%s.png", order.OrderNumber)
				uploadResult, err := s.storageSvc.UploadData(png, fileName, "image/png", FolderTickets)
				if err == nil {
					order.QRCodeURL = uploadResult.SecureURL
				}
			}
		}

		// Deduct quantity
		if err := txRepo.UpdateQuantity(ticketID, quantity); err != nil {
			return err
		}

		// Create the order record for tracing and revenue tracking
		return txRepo.CreateOrder(order)
	})
}

func (s *ticketService) PurchaseTicketGuest(order *models.TicketOrder) error {
	return s.repo.Transaction(func(txRepo repository.TicketRepository) error {
		ticket, err := txRepo.GetByID(order.TicketID.String())
		if err != nil {
			return errors.New("ticket not found")
		}

		if ticket.AvailableQuantity < order.Quantity {
			return errors.New("not enough tickets available")
		}

		order.TotalAmount = ticket.Price * float64(order.Quantity)

		if order.Category == "" {
			order.Category = ticket.Category
		}
		orderNumber, err := generateOrderNumber()
		if err != nil {
			return err
		}
		order.OrderNumber = orderNumber

		// Generate TICKET QR Code
		if s.storageSvc != nil {
			// Construct a validation URL so scanners recognize it as data/link
			qrContent := fmt.Sprintf("https://webuye-sportif-fc.devsinkenya.com/verify-ticket?order=%s", order.OrderNumber)
			png, err := qrcode.Encode(qrContent, qrcode.Medium, 256)
			if err != nil {
				return fmt.Errorf("failed to generate QR code: %w", err)
			}

			fileName := fmt.Sprintf("ticket_%s.png", order.OrderNumber)
			uploadResult, err := s.storageSvc.UploadData(png, fileName, "image/png", FolderTickets)
			if err == nil {
				order.QRCodeURL = uploadResult.SecureURL
			}
		}

		// Deduct quantity
		if err := txRepo.UpdateQuantity(order.TicketID.String(), order.Quantity); err != nil {
			return err
		}

		// Create the order record
		return txRepo.CreateOrder(order)
	})
}

func generateOrderNumber() (string, error) {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := 8
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return "TKT-" + string(result), nil
}

func (s *ticketService) CreateTicket(ticket *models.Ticket) error {
	return s.repo.Create(ticket)
}
