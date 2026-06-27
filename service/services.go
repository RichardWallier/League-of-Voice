package service

import entity "lov/repository"

type Services struct {
	UserService *UserService
	AuthService *AuthService
	TokenService *TokenService
	WebSocketService *WebSocketService
	SFUService *SFUService
}

func NewServices(e *entity.Entities) *Services{
	return &Services{
		NewUserService(e.UserEntity),
		NewAuthService(e.AuthEntity, e.UserEntity),
		NewTokenService(e.TokenEntity, e.UserEntity),
		NewWebSocketService(),
		NewSFUService(),
	}
}
