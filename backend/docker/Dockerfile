# ===================================================================
# Stage 1: Build Go binary
# ===================================================================
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk --no-cache add git ca-certificates tzdata

WORKDIR /app

# Atur GOPROXY dengan fallback 'direct' agar jika proxy golang timeout/terblokir,
# Go akan langsung mencoba download dari repository sumber (GitHub, dll.)
ENV GOPROXY=https://proxy.golang.org,direct
ENV GO111MODULE=on

# Copy dependency files lebih dulu untuk memanfaatkan layer caching Docker
COPY go.mod go.sum ./

# Download dependencies (ditambahkan -x atau retry jika koneksi lambat)
RUN go mod download

# Copy seluruh source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o /app/bin/server \
    ./cmd/server

# Build CLI installer
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" \
    -o /app/bin/installer \
    ./cmd/installer

# ===================================================================
# Stage 2: Minimal runtime image
# ===================================================================
FROM alpine:3.19

RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN addgroup -S hris && adduser -S hris -G hris

WORKDIR /app

# Copy binary and config from builder
COPY --from=builder /app/bin/server .
COPY --from=builder /app/bin/installer .
COPY --from=builder /app/config ./config
COPY --from=builder /app/migrations ./migrations

# Use non-root user
USER hris

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

EXPOSE 8080

ENTRYPOINT ["./server"]
CMD ["--config", "./config/config.yaml"]