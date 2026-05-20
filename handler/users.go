package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"lov/dto"
	"lov/entity"
	"lov/service"
	"net/http"
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
	fmt.Println("Received request to create user")
		var newUserRequest dto.CreateUserRequest
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
