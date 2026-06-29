# ---- Build stage ----
FROM golang:1.26-alpine AS builder

# ca-certificates so we can copy them into the scratch image (Postgres TLS, etc.)
RUN apk add --no-cache ca-certificates

WORKDIR /src

# Download deps first so this layer is cached unless go.mod/go.sum change
COPY go.mod go.sum ./
RUN go mod download

# Build a fully static binary (CGO disabled -> no libc needed at runtime)
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /server .

# ---- Runtime stage ----
FROM scratch

# TLS root certs (needed for outbound TLS, e.g. a managed Postgres)
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /server /server

# Documentational only (ignored under host networking on the VPS).
# HTTP/WS signaling:
EXPOSE 3000

ENTRYPOINT ["/server"]
