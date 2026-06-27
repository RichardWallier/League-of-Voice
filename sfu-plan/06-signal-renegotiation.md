# Step 6 — signalPeerConnections: the renegotiation engine

> [SFU plan](../SFU_PLAN.md) · Step 6 of 8 · prev: [05](./05-ontrack-fanout.md) · next: [07](./07-client-rewrite.md)

**Goal:** Every peer ends up sending every other peer's audio. This is what makes the call
actually work.

**Why:** When the track set changes (join/leave/`OnTrack`), each PeerConnection must be
updated and re-offered. Piece 3 of [How the SFU works](../SFU_PLAN.md#how-the-sfu-works--the-pieces).
Port `main.go:119-214` (drop the keyframe bits).

**Files:** `service/websocket.go`.

**Do:** implement `signalPeerConnections()`:

1. `s.listLock.Lock()`; `defer s.listLock.Unlock()`. (Example also dispatches a keyframe on
   unlock — **skip**, audio-only.)
2. Inner `attemptSync() (tryAgain bool)` loops over `s.peerConnections`; for each PC:
   - if `ConnectionState() == webrtc.PeerConnectionStateClosed`, remove it from the slice and
     `return true` (slice changed — restart).
   - build `existingSenders` from `GetSenders()`; for any sender whose track ID is **not** in
     `s.trackLocals`, `RemoveTrack(sender)` (`return true` on error).
   - also add each `GetReceivers()` track ID to `existingSenders` — this prevents a peer from
     **receiving its own audio** (no loopback).
   - for every `trackID` in `s.trackLocals` not in `existingSenders`,
     `AddTrack(s.trackLocals[trackID])`.
   - `offer, _ := pc.CreateOffer(nil)` → `pc.SetLocalDescription(offer)` →
     `b, _ := json.Marshal(offer)` → `WriteJSON(&websocketMessage{Event:"offer", Data:string(b)})`.
3. Outer loop: call `attemptSync()` up to **25** times; if it still returns `true` after 25,
   spawn `go func(){ time.Sleep(3*time.Second); s.signalPeerConnections() }()` and `return`.
   - **Why the retry/backoff:** you can't renegotiate a PeerConnection while a previous
     negotiation is in flight. Keep this loop — it's not optional.

**Verify:**
- `go build ./...` passes — this completes the server side.
- The "you can hear each other" milestone is confirmed in [step 8](./08-run-and-verify.md)
  once the client exists. If you already built the client, two tabs should now hear each
  other and a third should join the mix.

**Next:** [Step 7 — rewrite the client](./07-client-rewrite.md).
