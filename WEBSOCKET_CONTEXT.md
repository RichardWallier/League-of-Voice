# WebSocket + WebRTC Voice Call — Context

## Goal
Real-time voice call between browser clients using WebRTC, with the Go WebSocket server acting purely as a signaling relay.

## Architecture

```
Browser A  <──── WebRTC audio (peer-to-peer) ────>  Browser B
     │                                                    │
     └──── JSON signaling via WebSocket ────> Go server ──┘
                                             (broadcasts to all)
```

The Go server does NOT process audio. It only relays small JSON signaling messages (offer/answer/ICE candidates) to all connected clients. Each client filters messages using `from` and `to` fields in the JSON.

## Current file state

### `api/handler/websocket.go`
```go
package handler

import (
    "lov/service"
    "net/http"
)

type WebSocketHandler struct {
    websocketService *service.WebSocketService
}

func NewWebSocketHandler(websocketService *service.WebSocketService) *WebSocketHandler {
    go websocketService.BroadcastMessages()  // started once at app startup
    return &WebSocketHandler{
        websocketService: websocketService,
    }
}

func (h *WebSocketHandler) ConnectWebSocket(w http.ResponseWriter, r *http.Request) {
    cleanup, readMessages, err := h.websocketService.UpgradeConnection(w, r)
    if err != nil {
        http.Error(w, "failed to upgrade or connect WebSocket", http.StatusInternalServerError)
        return
    }
    defer cleanup()
    readMessages()
}
```

### `api/service/websocket.go`
```go
package service

import (
    "net/http"
    "sync"
    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []byte)
var mutex = &sync.Mutex{}

type WebSocketService struct{}

func NewWebSocketService() *WebSocketService {
    return &WebSocketService{}
}

func (s *WebSocketService) UpgradeConnection(w http.ResponseWriter, r *http.Request) (func(), func(), error) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        return nil, nil, err
    }
    mutex.Lock()
    clients[conn] = true
    mutex.Unlock()
    return func() {
        conn.Close()
    }, func() { readMessages(conn) }, nil
}

func readMessages(conn *websocket.Conn) {
    for {
        _, message, err := conn.ReadMessage()
        if err != nil {
            mutex.Lock()
            delete(clients, conn)
            mutex.Unlock()
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
        for client := range clients {
            err := client.WriteMessage(websocket.TextMessage, message)
            if err != nil {
                client.Close()
                delete(clients, client)
            }
        }
        mutex.Unlock()
    }
}
```

### `api/routes/websocket.go`
```go
package routes

import (
    "lov/handler"
    "github.com/go-chi/chi/v5"
)

func WebSocketRoutes(chi *chi.Mux, websocketHandler *handler.WebSocketHandler) {
    chi.Get("/ws", websocketHandler.ConnectWebSocket)
}
```

### `api/main.go`
```go
func main() {
    fmt.Println("Server starting...")
    ctx := context.Background()
    db := db.NewPostgresDB(ctx)
    defer db.Cleanup()
    entities := repository.NewEntities(db)
    services := service.NewServices(entities)
    handlers := handler.SetupHandlers(services)
    routes := routes.SetupRoutes(handlers)
    fmt.Println("Server running on :3000...")
    if err := http.ListenAndServe(":3000", routes); err != nil {
        fmt.Println(err.Error())
    }
}
```

### `api/test-client.html`
WebRTC voice call client. Each tab gets a random 6-char ID. Uses the WebSocket purely for signaling (offer/answer/ICE). Audio is peer-to-peer via WebRTC. Supports N-way calls (one RTCPeerConnection per peer). Uses Google's public STUN server.

Signaling message format:
- `{type: "join", from: id}` — broadcast on connect
- `{type: "offer", from, to, sdp}` — WebRTC offer
- `{type: "answer", from, to, sdp}` — WebRTC answer
- `{type: "ice", from, to, candidate}` — ICE candidate
- `{type: "leave", from}` — broadcast on disconnect

Glare resolution: higher string ID always initiates the offer.

## Known issues / things to investigate
- The server sends `websocket.TextMessage` in `BroadcastMessages` — signaling JSON is text so this is fine
- `broadcast` channel is unbuffered — if `BroadcastMessages` is slow, `readMessages` will block. Not a problem for signaling (tiny messages), but worth noting
- No authentication — anyone can connect to `/ws`
- No per-client filtering on the server — all messages are broadcast to all clients; filtering is done client-side via the `to` field

## How to run
1. `cd api && go run main.go`
2. Open `api/test-client.html` in two browser tabs
3. Click **Join Call** in both, grant mic permission
4. You should see the WebRTC handshake in the log and hear audio from the other tab (use headphones)

## Collaboration rules
- Do NOT edit or write `.go` files — explain changes and let the user apply them
- Non-Go files (HTML, SQL, JSON, configs) can be edited directly
