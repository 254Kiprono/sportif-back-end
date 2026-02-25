package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MatchHandler struct {
	service services.MatchService
}

func NewMatchHandler(service services.MatchService) *MatchHandler {
	return &MatchHandler{service}
}

func (h *MatchHandler) GetLineup(c *gin.Context) {
	fixtureID := c.Param("id")
	lineup, err := h.service.GetLineup(fixtureID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "lineup not found"})
		return
	}
	c.JSON(http.StatusOK, lineup)
}

func (h *MatchHandler) SaveLineup(c *gin.Context) {
	var lineup models.Lineup
	if err := c.ShouldBindJSON(&lineup); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.SaveLineup(&lineup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lineup)
}

func (h *MatchHandler) StartMatch(c *gin.Context) {
	fixtureID := c.Param("id")
	if err := h.service.StartMatch(fixtureID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match started, appearances updated"})
}

func (h *MatchHandler) EndMatch(c *gin.Context) {
	fixtureID := c.Param("id")
	if err := h.service.EndMatch(fixtureID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "match ended"})
}

func (h *MatchHandler) GetEvents(c *gin.Context) {
	fixtureID := c.Param("id")
	events, err := h.service.GetEvents(fixtureID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, events)
}

func (h *MatchHandler) LogEvent(c *gin.Context) {
	var event models.MatchEvent
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fixtureID, _ := uuid.Parse(c.Param("id"))
	event.FixtureID = fixtureID

	if err := h.service.LogEvent(&event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, event)
}
