package handler

import (
	"encoding/json"
	"errors"
	"lov/service"
	"lov/utils"
	"net/http"
)

type UserHandler struct {
	service *service.UserService
	tokenService *service.TokenService
}

func NewUserHandler(s *service.UserService, t *service.TokenService) *UserHandler {
	return &UserHandler{
		service: s,
		tokenService: t,
	}
}

func (h *UserHandler) GetAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := h.service.GetAllUsers(r.Context())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *UserHandler) GetCurrentUserHandler(w http.ResponseWriter, r *http.Request) {
	email, err := utils.ValidateJWTToken(r.Context(), r.Header.Get("Authorization"))
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
