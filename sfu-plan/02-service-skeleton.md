# Step 2 — Replace the relay with the SFUService skeleton

> [SFU plan](../SFU_PLAN.md) · Step 2 of 8 · prev: [01](./01-dependencies.md) · next: [03](./03-routes-timeout.md)

**Goal:** Swap the mesh relay for an `SFUService` struct with state + empty methods; the
project compiles and `/ws` still upgrades.

**Why:** You're changing the server's job from "broadcast bytes" to "be a WebRTC peer." This
step sets up the new shape (state on the struct, per decision #3) without media logic yet —
a clean compiling base. See [How the SFU works](../SFU_PLAN.md#how-the-sfu-works--the-pieces).

**Files:** `service/websocket.go` (rewrite), `service/services.go`, `handler/websocket.go`,
`handler/handler.go`.

**Do:**

1. In `service/websocket.go`, **delete the relay**: `clients`, `broadcast`, `mutex`,
   `readMessages`, `BroadcastMessages`. **Keep** `upgrader`. (Keep `frameType`/logging if you like.)

2. Add the wire type + state-bearing service (port from example globals `main.go:34-49`):
   ```go
   type websocketMessage struct {
       Event string `json:"event"`
       Data  string `json:"data"`
   }

   type peerConnectionState struct {
       peerConnection *webrtc.PeerConnection
       websocket      *threadSafeWriter
   }

   type SFUService struct {
       listLock        sync.RWMutex
       peerConnections []peerConnectionState
       trackLocals     map[string]*webrtc.TrackLocalStaticRTP
   }

   func NewSFUService() *SFUService {
       return &SFUService{trackLocals: map[string]*webrtc.TrackLocalStaticRTP{}}
   }
   ```

3. Add empty stubs so it compiles (you fill them in steps 4–6):
   ```go
   func (s *SFUService) Join(conn *websocket.Conn)                                    { /* step 4 */ }
   func (s *SFUService) addTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP   { /* step 5 */ return nil }
   func (s *SFUService) removeTrack(t *webrtc.TrackLocalStaticRTP)                    { /* step 5 */ }
   func (s *SFUService) signalPeerConnections()                                       { /* step 6 */ }
   ```

4. Add the `threadSafeWriter` helper (port `main.go:409-420`): embeds `*websocket.Conn` +
   `sync.Mutex`, exposes a locked `WriteJSON`. Needed in steps 4–6 (multiple goroutines write one socket).

5. `service/services.go`: rename `WebSocketService` → `SFUService`, `NewWebSocketService()` →
   `NewSFUService()`; update the `Services` field + `NewServices`.

6. `handler/websocket.go`: the handler should **upgrade the connection** and call
   `sfuService.Join(conn)`. Remove the old `cleanup/readMessages` dance, and remove the
   `go ...BroadcastMessages()` line in `NewWebSocketHandler` (no global broadcaster anymore).
   Keep your logs. (Rename to `SFUHandler` if you want, or keep the name.)

7. `handler/handler.go`: update the wiring to the renamed service/handler.

> The WS upgrade now lives in the **handler** because it needs the `*websocket.Conn` to pass
> to `Join`. Either move `upgrader.Upgrade` into the handler, or give the service an `Upgrade`
> method that returns the conn.

**Verify:**
- `go build ./...` passes.
- Run the server; in a browser console run `new WebSocket('ws://localhost:3000/ws')` — your
  upgrade/connect log should fire. Nothing else happens yet.

**Next:** [Step 3 — take `/ws` out of the 60s timeout](./03-routes-timeout.md).
