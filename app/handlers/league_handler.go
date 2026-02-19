package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
)

type LeagueHandler struct {
	service services.LeagueService
}

func NewLeagueHandler(service services.LeagueService) *LeagueHandler {
	return &LeagueHandler{service}
}

func (h *LeagueHandler) GetTable(c *gin.Context) {
	table, err := h.service.GetTable()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, table)
}

func (h *LeagueHandler) UpdateEntry(c *gin.Context) {
	id := c.Param("id")
	var entry models.LeagueTable
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateEntry(id, &entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *LeagueHandler) CreateEntry(c *gin.Context) {
	var entry models.LeagueTable
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateEntry(&entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, entry)
}
