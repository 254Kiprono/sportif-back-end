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
	query := `SELECT * FROM membership_plans WHERE deleted_at IS NULL`
	err := r.db.Raw(query).Scan(&plans).Error
	return plans, err
}

func (r *membershipRepository) GetPlanByID(id string) (*models.MembershipPlan, error) {
	var plan models.MembershipPlan
	query := `SELECT * FROM membership_plans WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, id).Scan(&plan).Error
	return &plan, err
}

func (r *membershipRepository) CreateOrder(order *models.MembershipOrder) error {
	query := `INSERT INTO membership_orders (id, created_at, updated_at, user_id, plan_id, status) VALUES (?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, order.ID, order.CreatedAt, order.UpdatedAt, order.UserID, order.PlanID, order.Status).Error
}

func (r *membershipRepository) CreatePlan(plan *models.MembershipPlan) error {
	query := `INSERT INTO membership_plans (id, created_at, updated_at, name, description, price, duration_months, benefits) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, plan.ID, plan.CreatedAt, plan.UpdatedAt, plan.Name, plan.Description, plan.Price,
		plan.DurationMonths, plan.Benefits).Error
}
