package handler

import (
	"log"
	"lov/service"
	"net/http"
)

type SFUHandler struct {
	websocketService	*service.WebSocketService
	sfuService				*service.SFUService
}

func NewSFUHandler(websocketService *service.WebSocketService, sfuService *service.SFUService) *SFUHandler {
	return &SFUHandler{
		websocketService: websocketService,
		sfuService:				sfuService,
	}
}

func (h *SFUHandler) Connect(w http.ResponseWriter, r *http.Request) {
	log.Printf("[ws] /ws hit: %s %s from %s", r.Method, r.URL.Path, r.RemoteAddr)

	unsafeConn, err := h.websocketService.UpgradeConnection(w, r)
	if err != nil {
		log.Printf("[ws] connection setup failed for %s: %v", r.RemoteAddr, err)
		http.Error(w, "failed to upgrade or connect WebSocket", http.StatusInternalServerError)
		return
	}

	h.sfuService.Join(unsafeConn)
}
