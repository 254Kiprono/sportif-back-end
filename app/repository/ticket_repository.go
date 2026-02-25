package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type TicketRepository interface {
	GetAll() ([]models.Ticket, error)
	GetByID(id string) (*models.Ticket, error)
	Create(ticket *models.Ticket) error
	Delete(id string) error
	GetAllOrders() ([]models.TicketOrder, error)
	CreateOrder(order *models.TicketOrder) error
	UpdateQuantity(id string, quantity int) error
	Transaction(fn func(repo TicketRepository) error) error
}

func (r *ticketRepository) Delete(id string) error {
	return r.db.Exec("UPDATE tickets SET deleted_at = NOW() WHERE id = ?", id).Error
}

func (r *ticketRepository) GetAllOrders() ([]models.TicketOrder, error) {
	var orders []models.TicketOrder
	err := r.db.Raw("SELECT * FROM ticket_orders WHERE deleted_at IS NULL ORDER BY created_at DESC").Scan(&orders).Error
	if err == nil {
		for i := range orders {
			r.db.Raw("SELECT * FROM tickets WHERE id = ? AND deleted_at IS NULL", orders[i].TicketID).Scan(&orders[i].Ticket)
			r.db.Raw("SELECT * FROM fixtures WHERE id = ? AND deleted_at IS NULL", orders[i].Ticket.FixtureID).Scan(&orders[i].Ticket.Fixture)
		}
	}
	return orders, err
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepository{db}
}

func (r *ticketRepository) GetAll() ([]models.Ticket, error) {
	var tickets []models.Ticket
	err := r.db.Raw("SELECT * FROM tickets WHERE deleted_at IS NULL").Scan(&tickets).Error
	if err == nil {
		for i := range tickets {
			r.db.Raw("SELECT * FROM fixtures WHERE id = ? AND deleted_at IS NULL", tickets[i].FixtureID).Scan(&tickets[i].Fixture)
		}
	}
	return tickets, err
}

func (r *ticketRepository) GetByID(id string) (*models.Ticket, error) {
	var ticket models.Ticket
	err := r.db.Raw("SELECT * FROM tickets WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&ticket).Error
	if err == nil {
		r.db.Raw("SELECT * FROM fixtures WHERE id = ? AND deleted_at IS NULL", ticket.FixtureID).Scan(&ticket.Fixture)
	}
	return &ticket, err
}

func (r *ticketRepository) Create(ticket *models.Ticket) error {
	return r.db.Exec("INSERT INTO tickets (id, created_at, updated_at, fixture_id, category, price, available_quantity) VALUES (?, ?, ?, ?, ?, ?, ?)",
		ticket.ID, ticket.CreatedAt, ticket.UpdatedAt, ticket.FixtureID, ticket.Category, ticket.Price, ticket.AvailableQuantity).Error
}

func (r *ticketRepository) CreateOrder(order *models.TicketOrder) error {
	return r.db.Exec("INSERT INTO ticket_orders (id, created_at, updated_at, ticket_id, full_name, mobile, email, category, order_number, quantity, total_amount, status, qr_code_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		order.ID, order.CreatedAt, order.UpdatedAt, order.TicketID, order.FullName, order.Mobile, order.Email, order.Category, order.OrderNumber, order.Quantity, order.TotalAmount, order.Status, order.QRCodeURL).Error
}

func (r *ticketRepository) UpdateQuantity(id string, quantity int) error {
	return r.db.Exec("UPDATE tickets SET available_quantity = available_quantity - ?, updated_at = NOW() WHERE id = ?", quantity, id).Error
}

func (r *ticketRepository) Transaction(fn func(repo TicketRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		return fn(NewTicketRepository(tx))
	})
}
