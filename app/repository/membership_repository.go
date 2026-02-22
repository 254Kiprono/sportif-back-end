package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type MembershipRepository interface {
	GetPlans() ([]models.MembershipPlan, error)
	GetPlanByID(id string) (*models.MembershipPlan, error)
	CreateOrder(order *models.MembershipOrder) error
	CreatePlan(plan *models.MembershipPlan) error
}

type membershipRepository struct {
	db *gorm.DB
}

func NewMembershipRepository(db *gorm.DB) MembershipRepository {
	return &membershipRepository{db}
}

func (r *membershipRepository) GetPlans() ([]models.MembershipPlan, error) {
	var plans []models.MembershipPlan
	err := r.db.Find(&plans).Error
	return plans, err
}

func (r *membershipRepository) GetPlanByID(id string) (*models.MembershipPlan, error) {
	var plan models.MembershipPlan
	err := r.db.First(&plan, "id = ?", id).Error
	return &plan, err
}

func (r *membershipRepository) CreateOrder(order *models.MembershipOrder) error {
	return r.db.Create(order).Error
}

func (r *membershipRepository) CreatePlan(plan *models.MembershipPlan) error {
	return r.db.Create(plan).Error
}
