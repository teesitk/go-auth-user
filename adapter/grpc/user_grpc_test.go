package grpcadapter_test

import (
	"context"
	"testing"

	grpcAdapter "go-auth-user/adapter/grpc"
	"go-auth-user/domain"
	pb "go-auth-user/proto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/mongo"
)

type MockUserRepoGRPC struct {
	mock.Mock
}

func (m *MockUserRepoGRPC) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepoGRPC) Find(id int) *domain.User {
	args := m.Called(id)
	return args.Get(0).(*domain.User)
}

func (m *MockUserRepoGRPC) FindByEmail(email string) *domain.User {
	args := m.Called(email)
	return args.Get(0).(*domain.User)
}

func (m *MockUserRepoGRPC) FindAll(page int, perPage int) []domain.User {
	args := m.Called(page, perPage)
	return args.Get(0).([]domain.User)
}

func (m *MockUserRepoGRPC) Update(id int, name string, email string) (*mongo.UpdateResult, error) {
	args := m.Called(id, name, email)
	fakeResult := &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}
	return fakeResult, args.Error(1)
}

func (m *MockUserRepoGRPC) Delete(id int) error {
	args := m.Called(id)
	return args.Error(1)
}

func (m *MockUserRepoGRPC) CountUsers() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestCreateUser(t *testing.T) {
	repo := new(MockUserRepoGRPC)
	server := grpcAdapter.NewUserServiceServer(repo)

	req := &pb.CreateUserRequest{Name: "Test", Email: "test@example.com", Password: "123"}
	repo.On("Create", mock.Anything).Return(nil)

	resp, err := server.CreateUser(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, "success", resp.Message)
}

func TestGetUser(t *testing.T) {
	repo := new(MockUserRepoGRPC)
	server := grpcAdapter.NewUserServiceServer(repo)

	repo.On("FindByEmail", "test@example.com").Return(&domain.User{Id: 1, Name: "Test", Email: "test@example.com"}, nil)
	resp, err := server.GetUser(context.Background(), &pb.GetUserRequest{Email: "test@example.com"})

	assert.NoError(t, err)
	assert.Equal(t, "Test", resp.Name)
	assert.Equal(t, "test@example.com", resp.Email)
}
