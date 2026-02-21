package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FanHandler struct {
	service services.FanService
}

func NewFanHandler(service services.FanService) *FanHandler {
	return &FanHandler{service}
}

func (h *FanHandler) GetAll(c *gin.Context) {
	fans, err := h.service.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fans)
}

func (h *FanHandler) Create(c *gin.Context) {
	var fan models.Fan
	if err := c.ShouldBindJSON(&fan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.Create(&fan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, fan)
}

func (h *FanHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var fan models.Fan
	if err := c.ShouldBindJSON(&fan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid fan ID format"})
		return
	}
	fan.ID = uID
	if err := h.service.Update(&fan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fan)
}

func (h *FanHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Fan deleted"})
}
