# Step 8 — Run & verify

> [SFU plan](../SFU_PLAN.md) · Step 8 of 8 · prev: [07](./07-client-rewrite.md)

**Goal:** Confirm the SFU works end-to-end.

**Files:** none.

**Do:**
1. `cd api && go run main.go` (or let Air rebuild).
2. Open `test-client.html` in two tabs; grant mic in both. **Use headphones** to avoid feedback.
3. Watch the server logs: per tab you should see connect → remote track arrived → offer sent →
   connection state `connected`.

**Pass criteria:**
- [ ] Two tabs hear each other.
- [ ] A third tab hears and is heard by both (N-way).
- [ ] **Kill the server mid-call → audio stops.** This proves media flows *through* the server
      (the SFU property — the opposite of the old mesh, where the call survived a server kill).
- [ ] A tab leaving stops its audio for the others, and the server logs `signalPeerConnections`
      re-syncing.

**If audio doesn't flow, check:**
- `s.trackLocals` is non-empty after someone talks (step 5).
- `signalPeerConnections` actually runs `AddTrack` (step 6) — and the no-loopback filter isn't
  over-filtering and skipping real tracks.
- offers reach the client (server log on send + client `setRemoteDescription` with no error).
- the client is sending `answer` back (network tab / logs).

**Done — that's Plan A.** Next: [Plan B — multiple rooms](../SFU_PLAN.md#plan-b--multiple-rooms-the-upgrade).
