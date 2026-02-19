package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
)

type FixtureHandler struct {
	service services.FixtureService
}

func NewFixtureHandler(service services.FixtureService) *FixtureHandler {
	return &FixtureHandler{service}
}

func (h *FixtureHandler) GetAll(c *gin.Context) {
	fixtures, err := h.service.GetFixtures()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fixtures)
}

func (h *FixtureHandler) Create(c *gin.Context) {
	var fixture models.Fixture
	if err := c.ShouldBindJSON(&fixture); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreateFixture(&fixture); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, fixture)
}

func (h *FixtureHandler) UpdateScore(c *gin.Context) {
	id := c.Param("id")
	var input struct {
		HomeScore int `json:"home_score"`
		AwayScore int `json:"away_score"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.UpdateScore(id, input.HomeScore, input.AwayScore); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Score updated successfully"})
}
