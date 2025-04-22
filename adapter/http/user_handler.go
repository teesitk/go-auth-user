package http

import (
	"encoding/json"
	"go-auth-user/domain"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type UserHandler struct {
	Service domain.UserService
}

type stdResp struct {
	Code    int
	Message string
}

type errResp struct {
	Code  int
	Error string
}

type listResp struct {
	stdResp
	Data []domain.User
}

type getResp struct {
	stdResp
	Data *domain.User
}

func NewUserHandler(service domain.UserService) *UserHandler {
	return &UserHandler{Service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	validate := validator.New(validator.WithRequiredStructEnabled())
	var req struct {
		Name     string `validate:"required"`
		Email    string `validate:"required,email"`
		Password string `validate:"required"`
	}
	var res getResp
	res.Code = http.StatusOK
	res.Message = "success"
	json.NewDecoder(r.Body).Decode(&req)
	if err := validate.Struct(req); err != nil {
		erJson, _ := json.MarshalIndent(&errResp{Code: http.StatusBadRequest, Error: err.Error()}, "", "    ")
		http.Error(w, string(erJson), http.StatusBadRequest)
		return
	}
	user, err := h.Service.CreateUser(req.Name, req.Email, req.Password)
	res.Data = user
	if err != nil {
		erJson, _ := json.MarshalIndent(&errResp{Code: http.StatusBadRequest, Error: err.Error()}, "", "    ")
		http.Error(w, string(erJson), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	var req struct{ Id int }
	var res getResp
	res.Code = http.StatusOK
	res.Message = "success"
	json.NewDecoder(r.Body).Decode(&req)
	user := h.Service.GetUser(req.Id)
	res.Data = user
	if user == nil {
		erJson, _ := json.MarshalIndent(&errResp{Code: http.StatusNotFound, Error: "user not found"}, "", "    ")
		http.Error(w, string(erJson), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) ListUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Page    int
		PerPage int
	}
	var res listResp
	res.Code = http.StatusOK
	res.Message = "success"

	json.NewDecoder(r.Body).Decode(&req)
	users := h.Service.ListUser(req.Page, req.PerPage)
	res.Data = users
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id    int
		Name  string
		Email string
	}
	json.NewDecoder(r.Body).Decode(&req)
	var res stdResp
	res.Code = http.StatusOK
	res.Message = "success"

	err := h.Service.UpdateUser(req.Id, req.Name, req.Email)
	if err != nil {
		erJson, _ := json.MarshalIndent(&errResp{Code: http.StatusBadRequest, Error: err.Error()}, "", "    ")
		http.Error(w, string(erJson), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Id int
	}
	json.NewDecoder(r.Body).Decode(&req)
	var res stdResp
	res.Code = http.StatusOK
	res.Message = "success"

	err := h.Service.DeleteUser(req.Id)
	if err != nil {
		erJson, _ := json.MarshalIndent(&errResp{Code: http.StatusBadRequest, Error: err.Error()}, "", "    ")
		http.Error(w, string(erJson), http.StatusBadRequest)
		return
	}
	json.NewEncoder(w).Encode(res)
}
