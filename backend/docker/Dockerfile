# ===================================================================
# Stage 1: Build Go binary
# ===================================================================
FROM golang:1.26-alpine AS builder

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
# Stage 2: Runtime image (Debian + LibreOffice Headless)
# ===================================================================
# Debian dipilih (bukan Alpine) agar LibreOffice (libreoffice-writer) dapat
# di-install via apt — dibutuhkan fitur Settings → Template Dokumen untuk
# konversi DOCX → PDF (preview & generate document).
FROM debian:bookworm-slim

# Install runtime deps + LibreOffice headless + font dasar.
# libreoffice-writer cukup untuk konversi .docx (tanpa seluruh suite);
# fonts-dejavu-core agar karakter/font dasar (termasuk aksen Indonesia)
# dirender dengan benar.
RUN apt-get update && apt-get install -y --no-install-recommends \
        ca-certificates \
        tzdata \
        curl \
        libreoffice-writer \
        fonts-dejavu-core \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user (dengan home dir — LibreOffice menulis user profile
# ke $HOME saat konversi pertama kali dijalankan)
RUN useradd --system --create-home --home-dir /home/hris --shell /usr/sbin/nologin hris
ENV HOME=/home/hris

WORKDIR /app

# Copy binary and config from builder.
#
# Migrasi TIDAK disalin: sejak direktori backend/migrations/ dihapus, seluruh
# file .sql ditanam ke dalam binary lewat //go:embed di
# internal/pkg/migrator/migrations/. Menyalinnya terpisah justru membuat versi
# migrasi di disk dan versi yang benar-benar dijalankan binary bisa berbeda
# tanpa ada yang memberi tahu.
COPY --from=builder /app/bin/server .
COPY --from=builder /app/bin/installer .
COPY --from=builder /app/config ./config

# Direktori upload (template .docx, previews/, generated_documents/) harus
# writable oleh user non-root
RUN mkdir -p /app/uploads && chown -R hris:hris /app /home/hris

# Use non-root user
USER hris

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/healthz || exit 1

EXPOSE 8080

ENTRYPOINT ["./server"]
CMD ["--config", "./config/config.yaml"]