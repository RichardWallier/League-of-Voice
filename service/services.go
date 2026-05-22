package service

import entity "lov/repository"

type Services struct {
	UserService *UserService
	AuthService *AuthService
	TokenService *TokenService
}

func NewServices(e *entity.Entities) *Services{
	userService := NewUserService(e.UserEntity)
	authService := NewAuthService(e.AuthEntity, e.UserEntity)
	tokenService := NewTokenService(e.TokenEntity, e.UserEntity)
	return &Services{
		userService,
		authService,
		tokenService,
	}
}
