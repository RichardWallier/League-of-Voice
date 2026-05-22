package handler

import (
	"encoding/json"
	"fmt"
	"lov/domain"
	"lov/dto"
	"lov/service"
	"lov/utils"
	"net/http"
)

type AuthHandler struct {
	authService *service.AuthService
	userService *service.UserService
	tokenService *service.TokenService
}

func NewAuthHandler(authService *service.AuthService, userService *service.UserService, tokenService *service.TokenService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		userService: userService,
		tokenService: tokenService,
	}
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
	token, err := h.authService.Login(r.Context(), loginRequest)
	if err != nil {
		http.Error(w, "invalid email or password", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Generated token: %s\n", token)
	w.Write([]byte(token))
}

func (h *AuthHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var registerRequest dto.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	if registerRequest.Email == "" || registerRequest.Username == "" || registerRequest.Password == "" {
		http.Error(w, "email, username and password are required", http.StatusBadRequest)
		return
	}

	userDomain := domain.User{
		Email:    registerRequest.Email,
		Username: registerRequest.Username,
		Password: registerRequest.Password,
	}
	user, err := h.userService.CreateUser(r.Context(), userDomain)
	if err != nil {
		http.Error(w, "failed to create user", http.StatusConflict)
		return
	}

	jwtToken, err := utils.GenerateJWTToken(r.Context(), user.Email)
	if err != nil {
		http.Error(w, "failed to login after registration", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Generated token after registration: %s\n", jwtToken)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": jwtToken})
}
