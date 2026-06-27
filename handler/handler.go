package handler

import "lov/service"

type Handlers struct {
	UserHandler *UserHandler
	AuthHandler *AuthHandler
	SFUHandler *SFUHandler
}

func SetupHandlers(s *service.Services) *Handlers {
	return &Handlers{
		UserHandler: NewUserHandler(s.UserService, s.TokenService),
		AuthHandler: NewAuthHandler(s.AuthService, s.UserService, s.TokenService),
		SFUHandler: NewSFUHandler(s.WebSocketService, s.SFUService),
	}
}
