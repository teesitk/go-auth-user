// Code generated by "mockgen"; DO NOT EDIT.
package domain

import (
	"github.com/stretchr/testify/mock"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return new(UserServiceMock)
}

func (_mock *UserServiceMock) CreateUser(name string, email string, password string) (*User, error) {
	_args := _mock.Called(name, email, password)
	var _a0 *User
	if _args.Get(0) != nil {
		if v, ok := _args.Get(0).(*User); ok {
			_a0 = v
		}
	}
	_a1 := _args.Error(1)
	return _a0, _a1
}
func (_mock *UserServiceMock) GetUser(id int) *User {
	_args := _mock.Called(id)
	var _a0 *User
	if _args.Get(0) != nil {
		if v, ok := _args.Get(0).(*User); ok {
			_a0 = v
		}
	}
	return _a0
}
func (_mock *UserServiceMock) ListUser(page int, perPage int) []User {
	_args := _mock.Called(page, perPage)
	var _a0 []User
	if _args.Get(0) != nil {
		if v, ok := _args.Get(0).([]User); ok {
			_a0 = v
		}
	}
	return _a0
}
func (_mock *UserServiceMock) UpdateUser(id int, name string, email string) error {
	_args := _mock.Called(id, name, email)
	_a0 := _args.Error(0)
	return _a0
}
func (_mock *UserServiceMock) DeleteUser(id int) error {
	_args := _mock.Called(id)
	_a0 := _args.Error(0)
	return _a0
}

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return new(UserRepositoryMock)
}

func (_mock *UserRepositoryMock) Create(user *User) error {
	_args := _mock.Called(user)
	_a0 := _args.Error(0)
	return _a0
}
func (_mock *UserRepositoryMock) Find(id int) *User {
	_args := _mock.Called(id)
	var _a0 *User
	if _args.Get(0) != nil {
		if v, ok := _args.Get(0).(*User); ok {
			_a0 = v
		}
	}
	return _a0
}
func (_mock *UserRepositoryMock) FindByEmail(email string) *User {
	_args := _mock.Called(email)
	var _a0 *User
	if _args.Get(0) != nil {
		if v, ok := _args.Get(0).(*User); ok {
			_a0 = v
		}
	}
	return _a0
}
func (_mock *UserRepositoryMock) FindAll(page int, perPage int) []User {
	_args := _mock.Called(page, perPage)
	var _a0 []User
	if _args.Get(0) != nil {
		if v, ok := _args.Get(0).([]User); ok {
			_a0 = v
		}
	}
	return _a0
}
func (_mock *UserRepositoryMock) Update(id int, name string, email string) (*mongo.UpdateResult, error) {
	_args := _mock.Called(id, name, email)
	var _a0 *mongo.UpdateResult
	if _args.Get(0) != nil {
		if v, ok := _args.Get(0).(*mongo.UpdateResult); ok {
			_a0 = v
		}
	}
	_a1 := _args.Error(1)
	return _a0, _a1
}
func (_mock *UserRepositoryMock) Delete(id int) error {
	_args := _mock.Called(id)
	_a0 := _args.Error(0)
	return _a0
}
func (_mock *UserRepositoryMock) CountUsers() (int64, error) {
	_args := _mock.Called()
	var _a0 int64
	if _args.Get(0) != nil {
		if v, ok := _args.Get(0).(int64); ok {
			_a0 = v
		}
	}
	_a1 := _args.Error(1)
	return _a0, _a1
}
