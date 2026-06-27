package routes

import (
	"lov/handler"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(r chi.Router, authHandler *handler.AuthHandler) {
	r.Post("/auth/login", authHandler.LoginHandler)
	r.Post("/auth/register", authHandler.RegisterHandler)
}
