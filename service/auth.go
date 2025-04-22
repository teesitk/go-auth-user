package service

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"go-auth-user/domain"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthServiceImpl struct {
	UserRepo  domain.UserRepository
	JWTSecret string
}

// Authenticate implements domain.AuthService.
func (a *AuthServiceImpl) Authenticate(user string, password string) (auth *domain.Auth, err error) {
	hash := md5.Sum([]byte(password))
	pass := hex.EncodeToString(hash[:])
	result := a.UserRepo.FindByEmail(user)
	if result.Email == user && result.Password == pass {
		//create token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub":   result.Id,
			"email": result.Email,
			"exp":   time.Now().Add(time.Hour * 24).Unix(),
		})
		signed, err := token.SignedString([]byte(a.JWTSecret))
		return &domain.Auth{Token: signed}, err
	}
	return &domain.Auth{}, errors.New("invalid user or password")
}

func (a *AuthServiceImpl) Signup(name string, email string, password string) error {
	hash := md5.Sum([]byte(password))
	user := &domain.User{
		Id:       0,
		Name:     name,
		Email:    email,
		Password: hex.EncodeToString(hash[:]),
	}
	err := a.UserRepo.Create(user)
	return err
}

func (a *AuthServiceImpl) ParseToken(tokenStr string) (*domain.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.JWTSecret), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return nil, errors.New("invalid token claims")
		}
		return a.UserRepo.FindByEmail(email), nil
	}

	return nil, err

}

func NewAuthService(repo domain.UserRepository, secret string) domain.AuthService {
	return &AuthServiceImpl{UserRepo: repo, JWTSecret: secret}
}
