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
	err := r.db.Raw("SELECT * FROM roles WHERE name = ? LIMIT 1", name).Scan(&role).Error
	if err == nil && role.ID != (uuid.UUID{}) {
		r.db.Raw("SELECT p.* FROM permissions p JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = ?", role.ID).Scan(&role.Permissions)
	}
	return &role, err
}

func (r *roleRepository) GetByID(id uuid.UUID) (*models.Role, error) {
	var role models.Role
	err := r.db.Raw("SELECT * FROM roles WHERE id = ? LIMIT 1", id).Scan(&role).Error
	if err == nil && role.ID != (uuid.UUID{}) {
		r.db.Raw("SELECT p.* FROM permissions p JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = ?", role.ID).Scan(&role.Permissions)
	}
	return &role, err
}

func (r *roleRepository) GetAll() ([]models.Role, error) {
	var roles []models.Role
	err := r.db.Raw("SELECT * FROM roles").Scan(&roles).Error
	if err == nil {
		for i := range roles {
			r.db.Raw("SELECT p.* FROM permissions p JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = ?", roles[i].ID).Scan(&roles[i].Permissions)
		}
	}
	return roles, err
}

func (r *roleRepository) Create(role *models.Role) error {
	err := r.db.Exec("INSERT INTO roles (id, created_at, updated_at, name, description) VALUES (?, ?, ?, ?, ?)",
		role.ID, role.CreatedAt, role.UpdatedAt, role.Name, role.Description).Error
	if err != nil {
		return err
	}
	for _, p := range role.Permissions {
		r.db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", role.ID, p.ID)
	}
	return nil
}

func (r *roleRepository) Update(role *models.Role) error {
	err := r.db.Exec("UPDATE roles SET name = ?, description = ?, updated_at = NOW() WHERE id = ?",
		role.Name, role.Description, role.ID).Error
	if err != nil {
		return err
	}
	// Sync permissions: delete existing and insert new
	r.db.Exec("DELETE FROM role_permissions WHERE role_id = ?", role.ID)
	for _, p := range role.Permissions {
		r.db.Exec("INSERT INTO role_permissions (role_id, permission_id) VALUES (?, ?)", role.ID, p.ID)
	}
	return nil
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	r.db.Exec("DELETE FROM role_permissions WHERE role_id = ?", id)
	return r.db.Exec("DELETE FROM roles WHERE id = ?", id).Error
}

func (r *roleRepository) GetPermissionsByRoleID(roleID uuid.UUID) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Raw("SELECT p.* FROM permissions p JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = ?", roleID).Scan(&permissions).Error
	return permissions, err
}
