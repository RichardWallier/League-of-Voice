package routes

import (
	"lov/handler"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(chi *chi.Mux, handlers *handler.Handlers) {
	chi.Post("/auth/login", handlers.AuthHandler.LoginHandler)
	// chi.Post("/auth/register", handlers.AuthHandler.RegisterHandler)
}
