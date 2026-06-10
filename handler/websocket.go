package handler

import (
	"log"
	"lov/service"
	"net/http"
)

type WebSocketHandler struct {
	websocketService *service.WebSocketService
}

func NewWebSocketHandler(websocketService *service.WebSocketService) *WebSocketHandler {
	log.Println("[ws] starting BroadcastMessages goroutine")
	go websocketService.BroadcastMessages()
	return &WebSocketHandler{
		websocketService: websocketService,
	}
}

func (h *WebSocketHandler) ConnectWebSocket(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ws] /ws hit: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	cleanup, readMessages, err := h.websocketService.UpgradeConnection(w, r)
	if err != nil {
		log.Printf("[ws] connection setup failed for %s: %v", r.RemoteAddr, err)
		http.Error(w, "failed to upgrade or connect WebSocket", http.StatusInternalServerError)
		return
	}
	defer cleanup()

	log.Printf("[ws] %s entering read loop", r.RemoteAddr)
	readMessages()
	log.Printf("[ws] %s read loop ended — handler returning", r.RemoteAddr)
}
