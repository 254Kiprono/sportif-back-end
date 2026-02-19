package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PlayerHandler struct {
	service services.PlayerService
}

func NewPlayerHandler(service services.PlayerService) *PlayerHandler {
	return &PlayerHandler{service}
}

func (h *PlayerHandler) GetAll(c *gin.Context) {
	players, err := h.service.GetPlayers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, players)
}

func (h *PlayerHandler) Create(c *gin.Context) {
	var player models.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreatePlayer(&player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, player)
}

func (h *PlayerHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var player models.Player
	if err := c.ShouldBindJSON(&player); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure external ID is consistent with URL param
	uID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid player ID format"})
		return
	}
	player.ID = uID

	if err := h.service.UpdatePlayer(&player); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error() + " for ID: " + id})
		return
	}
	c.JSON(http.StatusOK, player)
}

func (h *PlayerHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePlayer(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Player deleted"})
}
