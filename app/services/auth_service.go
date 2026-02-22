package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"webuye-sportif/app/config"
	"webuye-sportif/app/models"
	"webuye-sportif/app/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(fullName, username, email, phone, password string) error
	CreateUser(fullName, username, email, phone, password, roleName string) error
	Login(username, password string) (string, error)
	Logout(jti string) error
	GetAllUsers() ([]models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
	roleRepo repository.RoleRepository
	cfg      *config.Config
	rdb      *redis.Client
}

func NewAuthService(
	userRepo repository.UserRepository,
	roleRepo repository.RoleRepository,
	cfg *config.Config,
	rdb *redis.Client,
) AuthService {
	return &authService{userRepo: userRepo, roleRepo: roleRepo, cfg: cfg, rdb: rdb}
}

func (s *authService) Register(fullName, username, email, phone, password string) error {
	return s.CreateUser(fullName, username, email, phone, password, "user")
}

func (s *authService) CreateUser(fullName, username, email, phone, password, roleName string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	role, err := s.roleRepo.GetByName(roleName)
	if err != nil {
		return err
	}

	user := &models.User{
		FullName: fullName,
		Username: username,
		Email:    email,
		Phone:    phone,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	return s.userRepo.Create(user)
}

func (s *authService) Login(username, password string) (string, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	var permissions []string
	for _, p := range user.Role.Permissions {
		permissions = append(permissions, p.Name)
	}

	// Unique session ID embedded in the JWT
	jti := uuid.New().String()
	expTime := time.Hour * 24

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":         jti,
		"user_id":     user.ID,
		"role_id":     user.RoleID,
		"role_name":   user.Role.Name,
		"permissions": permissions,
		"exp":         time.Now().Add(expTime).Unix(),
	})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	// Store session in Redis using the injected client (whitelist approach)
	if s.rdb != nil {
		ctx := context.Background()
		key := fmt.Sprintf("session:%s", jti)
		if err = s.rdb.Set(ctx, key, user.ID.String(), expTime).Err(); err != nil {
			return "", errors.New("failed to create session")
		}
	}

	return tokenString, nil
}

func (s *authService) Logout(jti string) error {
	if s.rdb == nil {
		return nil
	}
	ctx := context.Background()
	return s.rdb.Del(ctx, fmt.Sprintf("session:%s", jti)).Err()
}

func (s *authService) GetAllUsers() ([]models.User, error) {
	return s.userRepo.GetAll()
}
