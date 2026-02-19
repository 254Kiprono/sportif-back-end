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
}

type authService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
}

func NewAuthService(userRepo repository.UserRepository, cfg *config.Config) AuthService {
	return &authService{userRepo, cfg}
}

func (s *authService) Register(fullName, username, email, phone, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &models.User{
		FullName: fullName,
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
		Role:     "user",
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

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.cfg.JWTSecret))
}
