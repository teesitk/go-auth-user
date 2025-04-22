package service_test

import (
	"errors"
	"go-auth-user/domain"
	"go-auth-user/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Save(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) CountUsers() (int64, error) {
	return 0, nil // not tested here
}

func TestSignup_Success(t *testing.T) {
	repo := new(MockUserRepo)
	svc := service.NewAuthService(repo, "secret")

	repo.On("FindByEmail", "test@example.com").Return(&domain.User{}, errors.New("not found"))
	repo.On("Save", mock.Anything).Return(nil)

	err := svc.Signup("Test", "test@example.com", "123")
	assert.NoError(t, err)
}

func TestLogin_Success(t *testing.T) {
	repo := new(MockUserRepo)
	svc := service.NewAuthService(repo, "secret")

	repo.On("FindByEmail", "test@example.com").Return(&domain.User{Email: "test@example.com", Password: "123"}, nil)

	token, err := svc.Login("test@example.com", "123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
