package handler

import (
	"encoding/json"
	"fmt"
	"lov/dto"
	"lov/service"
	"net/http"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(s *service.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var loginRequest dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if loginRequest.Email == "" || loginRequest.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}
	token, err := h.service.Login(r.Context(), loginRequest)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Generated token: %s\n", token)
	w.Write([]byte(token))
}
