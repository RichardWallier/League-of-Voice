package routes

import (
	"lov/handler"

	"github.com/go-chi/chi/v5"
)

func UsersRoutes(chi *chi.Mux, handlers *handler.Handlers) {
	chi.Get("/Users", handlers.UserHandler.GetAllUsersHandler)
	chi.Post("/Users", handlers.UserHandler.CreateUserHandler)

	chi.Get("/Users/me", handlers.UserHandler.GetCurrentUserHandler)
}
