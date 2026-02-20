package services

import (
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"

	"github.com/google/uuid"
)

type RoleService interface {
	GetRoleByName(name string) (*models.Role, error)
	GetRoleByID(id uuid.UUID) (*models.Role, error)
	CheckPermission(roleID uuid.UUID, permissionName string) (bool, error)
}

type roleService struct {
	repo repository.RoleRepository
}

func NewRoleService(repo repository.RoleRepository) RoleService {
	return &roleService{repo}
}

func (s *roleService) GetRoleByName(name string) (*models.Role, error) {
	return s.repo.GetByName(name)
}

func (s *roleService) GetRoleByID(id uuid.UUID) (*models.Role, error) {
	return s.repo.GetByID(id)
}

func (s *roleService) CheckPermission(roleID uuid.UUID, permissionName string) (bool, error) {
	role, err := s.repo.GetByID(roleID)
	if err != nil {
		return false, err
	}

	for _, p := range role.Permissions {
		if p.Name == permissionName {
			return true, nil
		}
	}

	return false, nil
}
