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
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (id, created_at, updated_at, full_name, username, email, phone, password, role_id) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, user.ID, user.CreatedAt, user.UpdatedAt, user.FullName, user.Username, user.Email, user.Phone, user.Password, user.RoleID).Error
}

func (r *userRepository) GetByUsername(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role.Permissions").Where("(username = ? OR email = ?) AND deleted_at IS NULL", identifier, identifier).First(&user).Error
	return &user, err
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role.Permissions").Where("id = ? AND deleted_at IS NULL", id).First(&user).Error
	return &user, err
}

func (r *userRepository) GetAll() ([]models.User, error) {
	var users []models.User
	err := r.db.Preload("Role.Permissions").Find(&users).Error
	return users, err
}
