package domain

type Auth struct {
	Token string
}

type AuthService interface {
	Signup(name string, email string, password string) error
	Authenticate(user string, password string) (*Auth, error)
	ParseToken(token string) (*User, error)
}

type AuthRepository interface {
}
