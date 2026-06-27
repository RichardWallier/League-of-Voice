package service

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketService struct {
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{}
}

func (s *WebSocketService) UpgradeConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	log.Printf("[ws] incoming connection from %s (origin=%q)", r.RemoteAddr, r.Header.Get("Origin"))
	fmt.Println("testing websocket connection upgrade")


	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade FAILED for %s: %v", r.RemoteAddr, err)
		return &websocket.Conn{}, err
	}

	return conn,  nil
}
