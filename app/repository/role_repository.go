package repository

import (
	"webuye-sportif/app/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RoleRepository interface {
	GetByName(name string) (*models.Role, error)
	GetByID(id uuid.UUID) (*models.Role, error)
	GetAll() ([]models.Role, error)
	Create(role *models.Role) error
	Update(role *models.Role) error
	Delete(id uuid.UUID) error
	GetPermissionsByRoleID(roleID uuid.UUID) ([]models.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db}
}

func (r *roleRepository) GetByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Where("name = ?", name).First(&role).Error
	return &role, err
}

func (r *roleRepository) GetByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").Where("id = ?", id).First(&role).Error
	return &role, err
}

func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Preload("Permissions").Find(&roles).Error
	return roles, err
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *roleRepository) GetPermissionsByRoleID(roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Model(&models.Role{BaseModel: models.BaseModel{ID: roleID}}).Association("Permissions").Find(&permissions)
	return permissions, err
}
