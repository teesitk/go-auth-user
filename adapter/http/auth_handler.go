package http

import (
	"encoding/json"
	"go-auth-user/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type AuthHandler struct {
	Service domain.AuthService
}

type stdAuthResp struct {
	Code    int
	Message string
}

type errAuthResp struct {
	Code  int
	Error string
}
type authResp struct {
	stdAuthResp
	Token string
}

func NewAuthHandler(service domain.AuthService) *AuthHandler {
	return &AuthHandler{Service: service}
}

func (h *AuthHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req struct {
		User     string
		Password string
	}
	var res authResp
	res.Code = http.StatusOK
	res.Message = "success"
	json.NewDecoder(r.Body).Decode(&req)
	auth, err := h.Service.Authenticate(req.User, req.Password)
	res.Token = auth.Token
	if err != nil {
		erJson, _ := json.MarshalIndent(&errResp{Code: http.StatusUnauthorized, Error: err.Error()}, "", "    ")
		http.Error(w, string(erJson), http.StatusUnauthorized)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {

	validate := validator.New(validator.WithRequiredStructEnabled())
	var req struct {
		Name     string `validate:"required"`
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}
	var res stdAuthResp
	res.Code = http.StatusOK
	res.Message = "success"
	json.NewDecoder(r.Body).Decode(&req)
	if err := validate.Struct(req); err != nil {
		erJson, _ := json.MarshalIndent(&errAuthResp{Code: http.StatusBadRequest, Error: err.Error()}, "", "    ")
		http.Error(w, string(erJson), http.StatusBadRequest)
		return
	}
	err := h.Service.Signup(req.Name, req.Email, req.Password)
	if err != nil {
		erJson, _ := json.MarshalIndent(&errAuthResp{Code: http.StatusBadRequest, Error: err.Error()}, "", "    ")
		http.Error(w, string(erJson), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(res)
}
