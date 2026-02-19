package handlers

import (
	"net/http"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
)

type DonationHandler struct {
	service services.DonationService
}

func NewDonationHandler(service services.DonationService) *DonationHandler {
	return &DonationHandler{service}
}

func (h *DonationHandler) Donate(c *gin.Context) {
	var input struct {
		Amount  float64 `json:"amount" binding:"required,gt=0"`
		Message string  `json:"message"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID *string
	if val, exists := c.Get("user_id"); exists {
		s := val.(string)
		userID = &s
	}

	donation, err := h.service.Donate(input.Amount, input.Message, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, donation)
}

func (h *DonationHandler) GetAll(c *gin.Context) {
	donations, err := h.service.GetDonations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, donations)
}
