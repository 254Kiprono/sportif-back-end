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
	err := r.db.Raw("SELECT * FROM membership_plans").Scan(&plans).Error
	return plans, err
}

func (r *membershipRepository) GetPlanByID(id string) (*models.MembershipPlan, error) {
	var plan models.MembershipPlan
	err := r.db.Raw("SELECT * FROM membership_plans WHERE id = ? LIMIT 1", id).Scan(&plan).Error
	return &plan, err
}

func (r *membershipRepository) CreateOrder(order *models.MembershipOrder) error {
	order.Initialize()
	return r.db.Exec("INSERT INTO membership_orders (id, created_at, updated_at, user_id, plan_id, status) VALUES (?, ?, ?, ?, ?, ?)",
		order.ID, order.CreatedAt, order.UpdatedAt, order.UserID, order.PlanID, order.Status).Error
}

func (r *membershipRepository) CreatePlan(plan *models.MembershipPlan) error {
	plan.Initialize()
	return r.db.Exec("INSERT INTO membership_plans (id, created_at, updated_at, name, description, price, duration_months, benefits) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		plan.ID, plan.CreatedAt, plan.UpdatedAt, plan.Name, plan.Description, plan.Price, plan.DurationMonths, plan.Benefits).Error
}
