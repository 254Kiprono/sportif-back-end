package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	GetByID(id string) (*models.User, error)
	GetAll() ([]models.User, error)
	Delete(id string) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Exec("INSERT INTO users (id, created_at, updated_at, full_name, username, email, phone, password, role_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		user.ID, user.CreatedAt, user.UpdatedAt, user.FullName, user.Username, user.Email, user.Phone, user.Password, user.RoleID).Error
}

func (r *userRepository) GetByUsername(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Raw("SELECT * FROM users WHERE (username = ? OR email = ? OR phone = ?) AND deleted_at IS NULL LIMIT 1", identifier, identifier, identifier).Scan(&user).Error
	if err == nil && user.ID != (models.BaseModel{}).ID {
		r.populateRole(&user)
	}
	return &user, err
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Raw("SELECT * FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1", id).Scan(&user).Error
	if err == nil && user.ID != (models.BaseModel{}).ID {
		r.populateRole(&user)
	}
	return &user, err
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Raw("SELECT * FROM users WHERE deleted_at IS NULL").Scan(&users).Error
	if err == nil {
		for i := range users {
			r.populateRole(&users[i])
		}
	}
	return users, err
}

func (r *userRepository) Delete(id string) error {
	return r.db.Exec("UPDATE users SET deleted_at = NOW() WHERE id = ?", id).Error
}

func (r *userRepository) populateRole(user *models.User) {
	r.db.Raw("SELECT * FROM roles WHERE id = ?", user.RoleID).Scan(&user.Role)
	r.db.Raw("SELECT p.* FROM permissions p JOIN role_permissions rp ON p.id = rp.permission_id WHERE rp.role_id = ?", user.RoleID).Scan(&user.Role.Permissions)
}
