package routes

import (
	"lov/handler"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(chi *chi.Mux, authHandler *handler.AuthHandler) {
	chi.Post("/auth/login", authHandler.LoginHandler)
	chi.Post("/auth/register", authHandler.RegisterHandler)
}
