package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"

	"github.com/google/uuid"
)

type MembershipService interface {
	GetPlans() ([]models.MembershipPlan, error)
	Subscribe(userID string, planID string) (*models.MembershipOrder, error)
	CreatePlan(plan *models.MembershipPlan) error
}

type membershipService struct {
	repo repository.MembershipRepository
}

func NewMembershipService(repo repository.MembershipRepository) MembershipService {
	return &membershipService{repo}
}

func (s *membershipService) GetPlans() ([]models.MembershipPlan, error) {
	return s.repo.GetPlans()
}

func (s *membershipService) Subscribe(userID string, planID string) (*models.MembershipOrder, error) {
	uID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	pID, err := uuid.Parse(planID)
	if err != nil {
		return nil, err
	}

	plan, err := s.repo.GetPlanByID(planID)
	if err != nil {
		return nil, err
	}

	order := &models.MembershipOrder{
		UserID: uID,
		PlanID: pID,
		Amount: plan.Price,
		Status: "pending",
	}

	err = s.repo.CreateOrder(order)
	return order, err
}

func (s *membershipService) CreatePlan(plan *models.MembershipPlan) error {
	return s.repo.CreatePlan(plan)
}
