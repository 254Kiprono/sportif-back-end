package services

import (
	"testing"
	"webuye-sportif/app/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepo) GetByID(id string) (*models.User, error) {
	args := m.Called(id)
	return args.Get(0).(*models.User), args.Error(1)
}

func TestAuthService_Register(t *testing.T) {
	mockRepo := new(MockUserRepo)
	service := NewAuthService(mockRepo, nil)

	mockRepo.On("Create", mock.Anything).Return(nil)

	err := service.Register("John Doe", "johndoe", "john@example.com", "0700000000", "password123")
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
