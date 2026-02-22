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
	return r.db.Create(user).Error
}

func (r *userRepository) GetByUsername(identifier string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Role.Permissions").Where("(username = ? OR email = ? OR phone = ?) AND deleted_at IS NULL", identifier, identifier, identifier).First(&user).Error
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

func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}
