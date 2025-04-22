package grpcadapter

import (
	"context"
	"errors"
	"net/http"

	"go-auth-user/domain"
	pb "go-auth-user/proto"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	UserRepo domain.UserRepository
}

func NewUserServiceServer(repo domain.UserRepository) *UserServiceServer {
	return &UserServiceServer{
		UserRepo: repo,
	}
}

func (s *UserServiceServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password, // In real-world: hash this!
	}

	if err := s.UserRepo.Create(user); err != nil {
		return &pb.CreateUserResponse{}, err
	}

	return &pb.CreateUserResponse{Code: http.StatusCreated, Message: "success"}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user := s.UserRepo.FindByEmail(req.Email)
	if user == nil {
		return nil, errors.New("user not found")
	}

	return &pb.GetUserResponse{
		Id:    "1",
		Name:  user.Name,
		Email: user.Email,
	}, nil
}
