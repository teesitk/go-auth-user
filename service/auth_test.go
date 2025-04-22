package service_test

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go-auth-user/domain"
	"go-auth-user/service"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepo) Find(id int) *domain.User {
	args := m.Called(id)
	return args.Get(0).(*domain.User)
}

func (m *MockUserRepo) FindByEmail(email string) *domain.User {
	args := m.Called(email)
	return args.Get(0).(*domain.User)
}

func (m *MockUserRepo) FindAll(page int, perPage int) []domain.User {
	args := m.Called(page, perPage)
	return args.Get(0).([]domain.User)
}

func (m *MockUserRepo) Update(id int, name string, email string) (*mongo.UpdateResult, error) {
	args := m.Called(id, name, email)
	fakeResult := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}
	return fakeResult, args.Error(1)
}

func (m *MockUserRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(1)
}

func (m *MockUserRepo) CountUsers() (int64, error) {
	return 0, nil // not tested here
}

func TestSignup_Success(t *testing.T) {
	repo := new(MockUserRepo)
	svc := service.NewAuthService(repo, "secret")

	repo.On("FindByEmail", "test@example.com").Return(&domain.User{}, errors.New("not found"))
	repo.On("Create", mock.Anything).Return(nil)

	err := svc.Signup("Test", "test@example.com", "123")
	assert.NoError(t, err)
}

func TestLogin_Success(t *testing.T) {
	repo := new(MockUserRepo)
	svc := service.NewAuthService(repo, "secret")

	hash := md5.Sum([]byte("123"))
	pass := hex.EncodeToString(hash[:])
	repo.On("FindByEmail", "test@example.com").Return(&domain.User{Email: "test@example.com", Password: pass}, nil)

	token, err := svc.Authenticate("test@example.com", "123")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
