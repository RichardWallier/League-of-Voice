package service

import "lov/entity"

type Services struct {
	UserService *UserService
	AuthService *AuthService
}

func NewServices(e *entity.Entities) *Services{
	return &Services{
		NewUserService(e.UserEntity),
		NewAuthService(e.AuthEntity),
	}
}
