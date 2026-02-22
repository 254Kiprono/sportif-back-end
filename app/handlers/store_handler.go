package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type StoreHandler struct {
	service services.StoreService
}

func NewStoreHandler(service services.StoreService) *StoreHandler {
	return &StoreHandler{service}
}

func (h *StoreHandler) GetJerseys(c *gin.Context) {
	jerseys, err := h.service.GetJerseys()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jerseys)
}

func (h *StoreHandler) GetOrders(c *gin.Context) {
	orders, err := h.service.GetOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *StoreHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateOrderStatus(id, input.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order status updated"})
}

func (h *StoreHandler) PlaceOrder(c *gin.Context) {
	userId, _ := c.Get("user_id")
	userIdStr := userId.(string)

	var input struct {
		Items []models.OrderItem `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, err := h.service.PlaceOrder(userIdStr, input.Items)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, order)
}

func (h *StoreHandler) Create(c *gin.Context) {
	var jersey models.Jersey
	if err := c.ShouldBindJSON(&jersey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateJersey(&jersey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, jersey)
}

func (h *StoreHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var jersey models.Jersey
	if err := c.ShouldBindJSON(&jersey); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid jersey ID format"})
		return
	}
	jersey.ID = uID

	if err := h.service.UpdateJersey(&jersey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, jersey)
}

func (h *StoreHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteJersey(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Jersey deleted"})
}
