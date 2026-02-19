package services

import (
	"errors"
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"
)

type TicketService interface {
	GetTickets() ([]models.Ticket, error)
	PurchaseTicket(ticketID string, userID string, quantity int) error
	CreateTicket(ticket *models.Ticket) error
}

type ticketService struct {
	repo repository.TicketRepository
}

func NewTicketService(repo repository.TicketRepository) TicketService {
	return &ticketService{repo}
}

func (s *ticketService) GetTickets() ([]models.Ticket, error) {
	return s.repo.GetAll()
}

func (s *ticketService) PurchaseTicket(ticketID string, userID string, quantity int) error {
	return s.repo.Transaction(func(txRepo repository.TicketRepository) error {
		ticket, err := txRepo.GetByID(ticketID)
		if err != nil {
			return errors.New("ticket not found")
		}

		if ticket.AvailableQuantity < quantity {
			return errors.New("not enough tickets available")
		}

		// Deduct quantity
		if err := txRepo.UpdateQuantity(ticketID, quantity); err != nil {
			return err
		}

		// In a real app, you'd also create a TicketOrder record here
		return nil
	})
}

func (s *ticketService) CreateTicket(ticket *models.Ticket) error {
	return s.repo.Create(ticket)
}
