package service

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var mutex = &sync.Mutex{}

type WebSocketService struct {
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{}
}

func (s *WebSocketService) UpgradeConnection(w http.ResponseWriter, r *http.Request) (func(), func(), error) {
	log.Printf("[ws] incoming connection from %s (origin=%q)", r.RemoteAddr, r.Header.Get("Origin"))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[ws] upgrade FAILED for %s: %v", r.RemoteAddr, err)
		return nil, nil, err
	}

	mutex.Lock()
	clients[conn] = true
	total := len(clients)
	mutex.Unlock()
	log.Printf("[ws] client connected: %s — total clients now %d", conn.RemoteAddr(), total)

	return func() {
		log.Printf("[ws] cleanup: closing connection %s", conn.RemoteAddr())
		conn.Close()
	}, func() { readMessages(conn) }, nil
}

func readMessages(conn *websocket.Conn) {
	// Read messages from the WebSocket connection
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			mutex.Lock()
			delete(clients, conn)
			remaining := len(clients)
			mutex.Unlock()
			log.Printf("[recv] read error from %s: %v — removed client, %d remaining", conn.RemoteAddr(), err, remaining)
			conn.Close()
			break
		}

		broadcast <- message
	}
}

func (s *WebSocketService) BroadcastMessages() {
	for {
		message := <-broadcast

		mutex.Lock()
		sent := 0
		for client := range clients {
			err := client.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				client.Close()
				delete(clients, client)
				continue
			}
			sent++
		}
		mutex.Unlock()
	}
}
