# SFU on LOV API — Implementation Plan (learning build)

Moving the LOV voice call from a **P2P mesh** (Go server = signaling relay) to an
**SFU** (Go server = WebRTC peer that forwards media). Built **from scratch with Pion**
to understand every piece. Sibling doc: `WEBSOCKET_CONTEXT.md` (current mesh).
Reference implementation to study: `example-webrtc-applications/sfu-ws/main.go`.

## Decisions (locked)

1. **Audio only.** Drops the entire video path *and* the RTCP keyframe logic.
2. **One global room first** (Plan A), then **multiple rooms** (Plan B). Learn the
   media mechanics with a single room before adding routing.
3. **State lives in the service struct**, not package globals (matches the DI style).
4. **Media flows through the server**, so restarting it drops every active call. Accepted.
5. **From scratch with Pion** — no LiveKit. This is a study project, not production.

---

## Mental model: mesh → SFU

| | Mesh (today) | SFU (this plan) |
|---|---|---|
| Server | Relays JSON bytes, never sees media | Is a WebRTC peer; receives + forwards audio |
| Connections per client | N−1 (one per peer) | 1 (just to the server) |
| Who offers | Clients offer each other | **Server** offers; client only answers |
| Kill server mid-call | Audio keeps flowing | Audio stops instantly |

---

## How the SFU works — the pieces

Walkthrough of one participant joining, mapped to `sfu-ws/main.go`. This is the core to
understand; every step in Plan A is a re-housing of this into your service struct.

**1. Connect & create the peer** (`main.go:237-276`)
- Upgrade the WebSocket.
- `webrtc.NewPeerConnection(...)` — the server's side of the call with this one client.
- Add a **recvonly audio transceiver**: "I expect to receive audio from you." The example
  loops over `[video, audio]` — **audio-only: keep just audio.**
- Append this PC to the room's `peerConnections` list.

**2. Wire the callbacks** (`main.go:279-350`)
- `OnICECandidate` → send each local candidate to the client as `{event:"candidate"}`.
- `OnConnectionStateChange` → `Failed`: close it. `Closed`: re-run `signalPeerConnections()`.
- `OnTrack` → **the fan-out.** Fires when the client's audio arrives. The server creates a
  `TrackLocalStaticRTP` (fan-out track), stores it in `trackLocals`, and loops reading RTP
  packets off the incoming track and `WriteRTP`s them into the fan-out track.

**3. `signalPeerConnections()` — the renegotiation engine** (`main.go:119-214`)
Runs whenever the track set changes (someone joins or leaves). For every PC in the room:
prune closed PCs, remove dead senders, **add every `trackLocals` entry the PC isn't sending
yet** (so each peer sends everyone else's audio), then create a fresh **offer** and send it.
Wrapped in a 25-try / 3s-backoff loop because you can't renegotiate mid-negotiation.

**4. The read loop** (`main.go:355-406`)
Reads from this client: `answer` → `SetRemoteDescription`; `candidate` → `AddICECandidate`.
The client never sends an offer.

**Why the server offers:** the server is the side whose track set changes (gains/loses
fan-out tracks as people come and go), so it drives renegotiation; the client just answers.

**The actual "selective forwarding":** A's packets enter via `OnTrack` → written to A's
`TrackLocalStaticRTP` → that local track was `AddTrack`'d onto B's and C's PeerConnections by
`signalPeerConnections` → so A's audio flows out to B and C.

---

## Plan A — single global room (audio-only), piece by piece

One room, no auth, all state on one `SFUService`, kept in **one file**
(`service/websocket.go`) to mirror `sfu-ws/main.go`. Build the steps in order — each ends
in a compiling state and tells you how to verify itself, so you can work through them solo.

| # | Step | Files | You'll have |
|---|------|-------|-------------|
| 1 | [Add Pion](./sfu-plan/01-dependencies.md) | `go.mod` | Pion installed, still compiles |
| 2 | [SFUService skeleton & wiring](./sfu-plan/02-service-skeleton.md) | `service/websocket.go`, `services.go`, `handler/*` | struct + empty methods, `/ws` upgrades |
| 3 | [Route off the 60s timeout](./sfu-plan/03-routes-timeout.md) | `routes/routes.go` | `/ws` no longer under `middleware.Timeout` |
| 4 | [Join + read loop](./sfu-plan/04-join-and-read-loop.md) | `service/websocket.go`, `handler/websocket.go` | a PeerConnection per client; answers/candidates handled |
| 5 | [OnTrack fan-out](./sfu-plan/05-ontrack-fanout.md) | `service/websocket.go` | incoming audio copied into shared fan-out tracks |
| 6 | [signalPeerConnections](./sfu-plan/06-signal-renegotiation.md) | `service/websocket.go` | renegotiation engine — everyone hears everyone |
| 7 | [Client rewrite](./sfu-plan/07-client-rewrite.md) | `test-client.html` | audio-only SFU client (server offers, client answers) |
| 8 | [Run & verify](./sfu-plan/08-run-and-verify.md) | — | audio across tabs; killing server drops the call |

Each step links back to [How the SFU works](#how-the-sfu-works--the-pieces) and to the
relevant `sfu-ws/main.go` lines. **Plan A is done when step 8 passes.**

---

## Plan B — multiple rooms (the upgrade)

The migration is clean because a global room is just *one* `Room`. Move the per-room state
off `SFUService` into a `Room`, and have `SFUService` own a registry of them.

```go
type Room struct {
    listLock        sync.RWMutex
    peerConnections []peerConnectionState
    trackLocals     map[string]*webrtc.TrackLocalStaticRTP
}

type SFUService struct {
    mu    sync.Mutex
    rooms map[string]*Room   // keyed by roomID
}
func (s *SFUService) room(id string) *Room              // get-or-create
func (s *SFUService) Join(roomID string, conn *websocket.Conn)
```

| File | Change |
|---|---|
| `service/room.go` | **new** — `Room` + `addTrack`/`removeTrack`/`signalPeerConnections` move here (operate on one room's lists). |
| `service/websocket.go` | `SFUService` holds `rooms map[string]*Room`; `Join` takes a `roomID`, looks up the room, delegates. Delete empty rooms when the last peer leaves. |
| `handler/websocket.go` | Read `roomID` from the URL, pass to `Join`. |
| `routes/websocket.go` | `GET /rooms/{roomID}/ws`. |
| `test-client.html` | Put the room id in the WS URL. |

Key idea: **the Plan A code becomes a single `Room`** — almost everything moves verbatim;
only the lookup/keying is new. (A detailed step file for this comes after Plan A is working.)

---

## Later / optional (only if you go beyond localhost)

- **Auth:** validate a JWT via your `TokenService` on connect. Browsers can't set headers on
  `new WebSocket()`, so people pass `?token=…` — but `middleware.Logger` + the raw-JSON
  logging would write the token to logs. Prefer **auth-as-first-message** or scrub it.
- **ICE/NAT:** add STUN/TURN + a public IP so remote clients can reach the server.
- **Graceful shutdown** in `main.go`: `http.Server` + `Shutdown`, close all PeerConnections.
- **Client reconnection** on socket drop.

---

## Gotchas

- **`middleware.Timeout(60s)` on `/ws`:** works today, but puts a bogus 60s deadline + a
  "write 504 on exit" on every long-lived socket. Scope the WS route out of it (step 3).
- **Thread-safe writes:** multiple goroutines write one socket (ICE callback + offer). Use the
  `threadSafeWriter` wrapper.
- **Renegotiation race:** keep the 25-try / 3s-backoff loop — you can't renegotiate mid-negotiation.
- **Candidate marshaling:** serialize outgoing candidates with `i.ToJSON()`, not the raw
  candidate, or `sdpMid` breaks.

---

## Next steps

- [ ] Plan A — work through `sfu-plan/01` → `08`
- [ ] Verify: audio across tabs, and killing the server drops the call
- [ ] Plan B — refactor into rooms
