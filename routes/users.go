package routes

import (
	"lov/handler"

	"github.com/go-chi/chi/v5"
)

func UsersRoutes(chi *chi.Mux, userHandler *handler.UserHandler) {
	chi.Get("/Users", userHandler.GetAllUsersHandler)

	chi.Get("/Users/me", userHandler.GetCurrentUserHandler)
}
