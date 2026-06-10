package routes

import (
	"lov/handler"

	"github.com/go-chi/chi/v5"
)

func WebSocketRoutes(chi *chi.Mux, websocketHandler *handler.WebSocketHandler) {
	chi.Get("/ws", websocketHandler.ConnectWebSocket)
}
