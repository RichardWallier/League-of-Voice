package routes

import (
	"lov/handler"

	"github.com/go-chi/chi/v5"
)

func SFURoutes(chi *chi.Mux, sfuhandler *handler.SFUHandler) {
	chi.Get("/ws", sfuhandler.Connect)
}
