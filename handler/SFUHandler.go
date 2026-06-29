package handler

import (
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

	unsafeConn, err := h.websocketService.UpgradeConnection(w, r)
	if err != nil {
		http.Error(w, "failed to upgrade or connect WebSocket", http.StatusInternalServerError)
		return
	}

	h.sfuService.Join(unsafeConn)
}
