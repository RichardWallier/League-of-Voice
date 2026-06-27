package routes

import (
	"lov/handler"

	"github.com/go-chi/chi/v5"
)

func UsersRoutes(r chi.Router, userHandler *handler.UserHandler) {
	r.Get("/Users", userHandler.GetAllUsersHandler)

	r.Get("/Users/me", userHandler.GetCurrentUserHandler)
}
