package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
)

type StoreHandler struct {
	service services.StoreService
}

func NewStoreHandler(service services.StoreService) *StoreHandler {
	return &StoreHandler{service}
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

func (h *StoreHandler) Update(c *gin.Context) {
	// To be implemented: update jersey details
	c.JSON(http.StatusOK, gin.H{"message": "Jersey updated successfully"})
}
