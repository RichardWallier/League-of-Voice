package handler

import (
	"encoding/json"
	"errors"
	"lov/dto"
	"lov/entity"
	"lov/service"
	"net/http"
	"strings"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{
		service: s,
	}
}

func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAllUsers(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUserRequest dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&newUserRequest); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if newUserRequest.Email == "" || newUserRequest.Username == "" || newUserRequest.Password == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	user, err := h.service.CreateUser(r.Context(), newUserRequest)
	if err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *UserHandler) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	email, err := service.ValidateJWTToken(r.Context(), strings.Trim(strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer"), " "))
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	user, err := h.service.GetUserByEmail(r.Context(), email)
	if err != nil {
		if errors.Is(err, service.ErrUserNotAuthenticated) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
