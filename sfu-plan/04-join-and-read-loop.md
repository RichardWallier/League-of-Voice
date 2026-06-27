# Step 4 — Join: PeerConnection, transceiver, callbacks, read loop

> [SFU plan](../SFU_PLAN.md) · Step 4 of 8 · prev: [03](./03-routes-timeout.md) · next: [05](./05-ontrack-fanout.md)

**Goal:** Each connecting client gets a server-side `PeerConnection`, and the server handles
the client's `answer` and ICE `candidate` messages. (Media fan-out is step 5.)

**Why:** This is the per-connection setup — pieces 1, 2 (minus OnTrack), and 4 of
[How the SFU works](../SFU_PLAN.md#how-the-sfu-works--the-pieces). Port `main.go:237-406`, audio-only.

**Files:** `service/websocket.go` (`Join`), `handler/websocket.go`.

**Do:** fill `Join(conn *websocket.Conn)`:

1. Wrap the conn: `c := &threadSafeWriter{conn, sync.Mutex{}}`; `defer c.Close()`.
2. `peerConnection, err := webrtc.NewPeerConnection(webrtc.Configuration{})`; `defer peerConnection.Close()`.
3. **Audio-only transceiver** (example loops video+audio at `main.go:263-271` — keep audio only):
   ```go
   if _, err := peerConnection.AddTransceiverFromKind(
       webrtc.RTPCodecTypeAudio,
       webrtc.RTPTransceiverInit{Direction: webrtc.RTPTransceiverDirectionRecvonly},
   ); err != nil { /* log + return */ }
   ```
4. Register the PC under the lock: append `peerConnectionState{peerConnection, c}` to `s.peerConnections`.
5. `OnICECandidate` (`main.go:279-300`): on non-nil candidate,
   `b, _ := json.Marshal(i.ToJSON())` then `c.WriteJSON(&websocketMessage{Event:"candidate", Data: string(b)})`.
   - **Gotcha:** use `i.ToJSON()`, not `i` — marshaling the raw candidate breaks `sdpMid`.
6. `OnConnectionStateChange` (`main.go:303-315`): `Failed` → `peerConnection.Close()`;
   `Closed` → `s.signalPeerConnections()`.
7. `OnTrack` → leave a `// TODO step 5` stub for now.
8. Call `s.signalPeerConnections()` once at the end of setup (no-op until step 6).
9. **Read loop** (`main.go:355-406`): loop `conn.ReadMessage()`, `json.Unmarshal` into a
   `websocketMessage`, switch on `Event`:
   - `"answer"`: `json.Unmarshal([]byte(msg.Data), &sd)` into a `webrtc.SessionDescription`,
     then `peerConnection.SetRemoteDescription(sd)`.
   - `"candidate"`: `json.Unmarshal([]byte(msg.Data), &ci)` into a `webrtc.ICECandidateInit`,
     then `peerConnection.AddICECandidate(ci)`.
   - **Double decode:** the outer frame is JSON, and `Data` is itself a JSON string.

In `handler/websocket.go`: upgrade, then call `sfuService.Join(conn)`.

> **Tip:** if you'd rather watch progress as you go, do [Step 7 (client)](./07-client-rewrite.md)
> now — then steps 5–6 will show audio appearing live. Otherwise verify by logs below.

**Verify:**
- `go build ./...` passes.
- Connect a client (the step-7 client, or a manual `RTCPeerConnection` that answers). Server
  logs should show ICE candidates flowing both ways and `OnConnectionStateChange` reaching
  `connected`. No audio yet (no fan-out).

**Next:** [Step 5 — OnTrack fan-out](./05-ontrack-fanout.md).
