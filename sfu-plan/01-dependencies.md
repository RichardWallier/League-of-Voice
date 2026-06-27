# Step 1 — Add Pion

> [SFU plan](../SFU_PLAN.md) · Step 1 of 8 · next: [02 — service skeleton](./02-service-skeleton.md)

**Goal:** Pion installed, project still compiles. No source changes yet.

**Why:** The SFU needs a real WebRTC stack on the server. Today you only have
`gorilla/websocket`. Pion (`github.com/pion/webrtc/v4`) gives you `PeerConnection`, tracks,
and RTP/RTCP — the same library the `sfu-ws` example uses.

**Files:** `go.mod`, `go.sum` (auto).

**Do:**
```sh
cd api
go get github.com/pion/webrtc/v4
go mod tidy
```

**Verify:**
- `go.mod` lists `github.com/pion/webrtc/v4` under `require`, and `gorilla/websocket` is no
  longer marked `// indirect`.
- `go build ./...` passes (nothing changed yet).

**Next:** [Step 2 — replace the relay with the SFUService skeleton](./02-service-skeleton.md).
