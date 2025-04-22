package domain

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	Id        int
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

type UserService interface {
	CreateUser(name string, email string, password string) (*User, error)
	GetUser(id int) *User
	ListUser(page int, perPage int) []User
	UpdateUser(id int, name string, email string) error
	DeleteUser(id int) error
}

type UserRepository interface {
	Create(user *User) error
	Find(id int) *User
	FindByEmail(email string) *User
	FindAll(page int, perPage int) []User
	Update(id int, name string, email string) (*mongo.UpdateResult, error)
	Delete(id int) error
	CountUsers() (int64, error)
}
