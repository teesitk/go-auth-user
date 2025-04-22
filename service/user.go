package service

import (
	"crypto/md5"
	"encoding/hex"
	"go-auth-user/domain"
)

type UserServiceImpl struct {
	Repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &UserServiceImpl{Repo: repo}
}

// CreateUser implements domain.UserService.
func (u *UserServiceImpl) CreateUser(name string, email string, password string) (*domain.User, error) {
	hash := md5.Sum([]byte(password))
	user := &domain.User{
		Id:       0,
		Name:     name,
		Email:    email,
		Password: hex.EncodeToString(hash[:]),
	}
	err := u.Repo.Create(user)
	return user, err
}

// DeleteUser implements domain.UserService.
func (u *UserServiceImpl) DeleteUser(id int) error {
	err := u.Repo.Delete(id)
	return err
}

// GetUser implements domain.UserService.
func (u *UserServiceImpl) GetUser(id int) *domain.User {
	user := u.Repo.Find(id)
	return user
}

// ListUser implements domain.UserService.
func (u *UserServiceImpl) ListUser(page int, perPage int) []domain.User {
	users := u.Repo.FindAll(page, perPage)
	return users
}

// UpdateUser implements domain.UserService.
func (u *UserServiceImpl) UpdateUser(id int, name string, email string) error {
	_, err := u.Repo.Update(id, name, email)
	return err
}
