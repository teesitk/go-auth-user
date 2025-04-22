package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	httpAdapter "go-auth-user/adapter/http"
	"go-auth-user/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAuthService struct {
	mock.Mock
}

func (m *mockAuthService) Signup(name, email, password string) error {
	return nil
}

func (m *mockAuthService) Authenticate(user, password string) (*domain.Auth, error) {
	return nil, nil
}
func (m *mockAuthService) ParseToken(token string) (*domain.User, error) {
	return nil, nil
}

func TestSignupHandler_Success(t *testing.T) {
	w := httptest.NewRecorder()
	body := map[string]string{"name": "Test", "email": "test@example.com", "password": "123"}
	b, _ := json.Marshal(body)
	r := httptest.NewRequest("POST", "/signup", bytes.NewBuffer(b))
	r.Header.Set("Content-Type", "application/json")

	h := httpAdapter.NewAuthHandler(&mockAuthService{})
	h.Signup(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
}
