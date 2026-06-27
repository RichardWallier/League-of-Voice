# Step 3 — Take `/ws` out of the 60s timeout

> [SFU plan](../SFU_PLAN.md) · Step 3 of 8 · prev: [02](./02-service-skeleton.md) · next: [04](./04-join-and-read-loop.md)

**Goal:** `/ws` no longer runs under `middleware.Timeout(60s)`.

**Why:** `routes/routes.go` applies `middleware.Timeout(60 * time.Second)` to every route. For
a long-lived WebSocket that's wrong — it puts a 60s deadline on the request context and tries
to write a `504` when the handler finally returns. It "works" today only because the hijacked
socket ignores the context, but it's a landmine once SFU code reads `r.Context()`. Keep the
timeout for REST; drop it for the socket. See [Gotchas](../SFU_PLAN.md#gotchas).

**Files:** `routes/routes.go`.

**Do:** Split the router with chi groups so the timeout wraps only REST:
```go
func SetupRoutes(handlers *handler.Handlers) *chi.Mux {
    r := chi.NewRouter()
    r.Use(middleware.Logger)

    // REST: timeout applies
    r.Group(func(r chi.Router) {
        r.Use(middleware.Timeout(60 * time.Second))
        UsersRoutes(r, handlers.UserHandler)
        AuthRoutes(r, handlers.AuthHandler)
    })

    // WebSocket: no timeout
    WebSocketRoutes(r, handlers.WebSocketHandler)

    fmt.Println("routers registered")
    return r
}
```
(`WebSocketRoutes` still registers `GET /ws` on the router it's handed.)

**Verify:**
- `go build ./...` passes; `/ws` still connects; REST routes still work.
- Optional: hold a `/ws` connection open > 60s — it stays open, no `504`-on-close artifact.

**Next:** [Step 4 — Join + read loop](./04-join-and-read-loop.md).
