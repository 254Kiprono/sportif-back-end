package repository

import (
	"webuye-sportif/app/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByUsername(username string) (*models.User, error)
	GetByID(id string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) Create(user *models.User) error {
	query := `INSERT INTO users (id, created_at, updated_at, full_name, username, email, phone, password, role) 
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	return r.db.Exec(query, user.ID, user.CreatedAt, user.UpdatedAt, user.FullName, user.Username, user.Email, user.Phone, user.Password, user.Role).Error
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE username = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, username).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.Username == "" {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}

func (r *userRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = ? AND deleted_at IS NULL LIMIT 1`
	err := r.db.Raw(query, id).Scan(&user).Error
	if err != nil {
		return nil, err
	}
	if user.ID.String() == "00000000-0000-0000-0000-000000000000" {
		return nil, gorm.ErrRecordNotFound
	}
	return &user, nil
}
