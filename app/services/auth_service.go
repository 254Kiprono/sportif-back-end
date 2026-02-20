package services

import (
	"errors"
	"time"

	"webuye-sportif/app/config"
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(fullName, username, email, phone, password string) error
	Login(username, password string) (string, error)
	GetAllUsers() ([]models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, roleRepo repository.RoleRepository, cfg *config.Config) AuthService {
	return &authService{userRepo, roleRepo, cfg}
}

func (s *authService) Register(fullName, username, email, phone, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	role, err := s.roleRepo.GetByName("user")
	if err != nil {
		return err
	}

	user := &models.User{
		BaseModel: models.BaseModel{},
		FullName:  fullName,
		Username:  username,
		Email:     email,
		Phone:     phone,
		Password:  string(hashedPassword),
		RoleID:    role.ID,
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	var permissions []string
	for _, p := range user.Role.Permissions {
		permissions = append(permissions, p.Name)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":     user.ID,
		"role_id":     user.RoleID,
		"role_name":   user.Role.Name,
		"permissions": permissions,
		"exp":         time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.cfg.JWTSecret))
}

func (s *authService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}
