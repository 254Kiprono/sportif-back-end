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
	query := `SELECT * FROM tickets WHERE deleted_at IS NULL`
	err := r.db.Raw(query).Scan(&tickets).Error
	return tickets, err
}

func (r *ticketRepository) GetByID(id string) (*models.Ticket, error) {
	var ticket models.Ticket
	query := `SELECT * FROM tickets WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, id).Scan(&ticket).Error
	return &ticket, err
}

func (r *ticketRepository) Create(ticket *models.Ticket) error {
	query := `INSERT INTO tickets (id, created_at, updated_at, fixture_id, category, price, available_quantity) 
	          VALUES (?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, ticket.ID, ticket.CreatedAt, ticket.UpdatedAt, ticket.FixtureID, ticket.Category,
		ticket.Price, ticket.AvailableQuantity).Error
}

func (r *ticketRepository) UpdateQuantity(id string, quantity int) error {
	query := `UPDATE tickets SET available_quantity = available_quantity - ? WHERE id = ?`
	return r.db.Exec(query, quantity, id).Error
}

func (r *ticketRepository) Transaction(fn func(repo TicketRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewTicketRepository(tx))
	})
}
