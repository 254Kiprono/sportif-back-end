package handlers

import (
	"net/http"
	"webuye-sportif/app/models"
	"webuye-sportif/app/services"

	"github.com/gin-gonic/gin"
)

type MembershipHandler struct {
	service services.MembershipService
}

func NewMembershipHandler(service services.MembershipService) *MembershipHandler {
	return &MembershipHandler{service}
}

func (h *MembershipHandler) GetPlans(c *gin.Context) {
	plans, err := h.service.GetPlans()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plans)
}

func (h *MembershipHandler) Subscribe(c *gin.Context) {
	userId, _ := c.Get("user_id")
	var input struct {
		PlanID string `json:"plan_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.service.Subscribe(userId.(string), input.PlanID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func (h *MembershipHandler) CreatePlan(c *gin.Context) {
	var plan models.MembershipPlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.CreatePlan(&plan); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plan)
}
