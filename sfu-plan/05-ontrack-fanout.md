# Step 5 — OnTrack: the fan-out

> [SFU plan](../SFU_PLAN.md) · Step 5 of 8 · prev: [04](./04-join-and-read-loop.md) · next: [06](./06-signal-renegotiation.md)

**Goal:** Audio arriving from one client is copied into a shared "local" track that can be
sent to the others. This is the actual *selective forwarding*.

**Why:** Pieces from [How the SFU works](../SFU_PLAN.md#how-the-sfu-works--the-pieces) step 2
(`OnTrack`) + the `addTrack`/`removeTrack` helpers. Port `main.go:88-116` and `main.go:317-346`.

**Files:** `service/websocket.go`.

**Do:**

1. `addTrack(t *webrtc.TrackRemote) *webrtc.TrackLocalStaticRTP` (`main.go:88-105`):
   - `s.listLock.Lock()`, and `defer` a func that **unlocks then calls `s.signalPeerConnections()`**
     (adding a track must trigger renegotiation),
   - `trackLocal, _ := webrtc.NewTrackLocalStaticRTP(t.Codec().RTPCodecCapability, t.ID(), t.StreamID())`,
   - `s.trackLocals[t.ID()] = trackLocal`; `return trackLocal`.
2. `removeTrack(t *webrtc.TrackLocalStaticRTP)` (`main.go:108-116`): same lock + defer
   (unlock then signal); `delete(s.trackLocals, t.ID())`.
3. Fill `OnTrack` in `Join` (`main.go:317-346`):
   ```go
   peerConnection.OnTrack(func(t *webrtc.TrackRemote, _ *webrtc.RTPReceiver) {
       trackLocal := s.addTrack(t)
       defer s.removeTrack(trackLocal)

       buf := make([]byte, 1500)
       rtpPkt := &rtp.Packet{}
       for {
           i, _, err := t.Read(buf)
           if err != nil { return }
           if err = rtpPkt.Unmarshal(buf[:i]); err != nil { return }
           rtpPkt.Extension = false
           rtpPkt.Extensions = nil
           if err = trackLocal.WriteRTP(rtpPkt); err != nil { return }
       }
   })
   ```
   - **Audio-only:** no `kind` filtering, and **no keyframe/PLI** (`dispatchKeyFrame` is
     video-only — skip it entirely).

**Verify:**
- `go build ./...` passes.
- With a client connected and talking, server logs should show a remote track arriving, and
  `len(s.trackLocals)` equals the number of talking clients.
- Still likely no audio until step 6 wires these tracks onto the *other* peers.

**Next:** [Step 6 — signalPeerConnections](./06-signal-renegotiation.md).
