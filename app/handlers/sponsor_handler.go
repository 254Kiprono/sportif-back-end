package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SponsorHandler struct {
	service services.SponsorService
}

func NewSponsorHandler(service services.SponsorService) *SponsorHandler {
	return &SponsorHandler{service}
}

func (h *SponsorHandler) GetAll(c *gin.Context) {
	sponsors, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sponsors)
}

func (h *SponsorHandler) Create(c *gin.Context) {
	var sponsor models.Sponsor
	if err := c.ShouldBindJSON(&sponsor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&sponsor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, sponsor)
}

func (h *SponsorHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var sponsor models.Sponsor
	if err := c.ShouldBindJSON(&sponsor); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sponsor ID format"})
		return
	}
	sponsor.ID = uID
	if err := h.service.Update(&sponsor); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sponsor)
}

func (h *SponsorHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Sponsor deleted"})
}
