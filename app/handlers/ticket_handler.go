package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TicketHandler struct {
	service services.TicketService
}

func NewTicketHandler(service services.TicketService) *TicketHandler {
	return &TicketHandler{service}
}

func (h *TicketHandler) GetAll(c *gin.Context) {
	tickets, err := h.service.GetTickets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

func (h *TicketHandler) Purchase(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var input struct {
		TicketID string `json:"ticket_id" binding:"required"`
		Quantity int    `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.PurchaseTicket(input.TicketID, userId.(string), input.Quantity); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket purchased successfully"})
}

func (h *TicketHandler) PurchaseGuest(c *gin.Context) {
	var input struct {
		TicketID string `json:"ticket_id" binding:"required"`
		FullName string `json:"full_name" binding:"required"`
		Mobile   string `json:"mobile" binding:"required"`
		Email    string `json:"email"`
		Category string `json:"category"`
		Quantity int    `json:"quantity" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse UUID
	tID, err := uuid.Parse(input.TicketID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket_id"})
		return
	}

	order := models.TicketOrder{
		TicketID: tID,
		FullName: input.FullName,
		Mobile:   input.Mobile,
		Email:    input.Email,
		Category: input.Category,
		Quantity: input.Quantity,
		Status:   "pending",
	}

	if err := h.service.PurchaseTicketGuest(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Ticket order placed successfully",
		"order":   order,
	})
}

func (h *TicketHandler) Create(c *gin.Context) {
	var ticket models.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateTicket(&ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, ticket)
}

func (h *TicketHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTicket(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Ticket deleted"})
}

func (h *TicketHandler) GetOrders(c *gin.Context) {
	orders, err := h.service.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
