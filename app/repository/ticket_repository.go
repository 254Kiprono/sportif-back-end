package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type TicketRepository interface {
	GetAll() ([]models.Ticket, error)
	GetByID(id string) (*models.Ticket, error)
	Create(ticket *models.Ticket) error
	UpdateQuantity(id string, quantity int) error
	Transaction(fn func(repo TicketRepository) error) error
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}

func (r *ticketRepository) GetAll() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Find(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetByID(id string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.First(&ticket, "id = ?", id).Error
	return &ticket, err
}

func (r *ticketRepository) Create(ticket *models.Ticket) error {
	return r.db.Create(ticket).Error
}

func (r *ticketRepository) UpdateQuantity(id string, quantity int) error {
	return r.db.Model(&models.Ticket{}).Where("id = ?", id).Update("available_quantity", gorm.Expr("available_quantity - ?", quantity)).Error
}

func (r *ticketRepository) Transaction(fn func(repo TicketRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewTicketRepository(tx))
	})
}
