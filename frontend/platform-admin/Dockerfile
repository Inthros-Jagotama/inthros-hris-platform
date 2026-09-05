# ===================================================================
# Stage 1: Build bundle Vite
# ===================================================================
FROM node:22-alpine AS builder

WORKDIR /app

COPY package.json package-lock.json ./
RUN npm ci

COPY . .
RUN npm run build

# ===================================================================
# Stage 2: Serve static via nginx
# ===================================================================
FROM nginx:1.27-alpine

# Sama seperti app tenant: src/services/api.js memakai `baseURL: ''`, jadi
# bundle ini harus diserve satu origin dengan backend. nginx.conf mem-proxy
# /api ke service `api`.
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80


# 127.0.0.1, bukan localhost: nginx.conf custom di sini bikin entrypoint
# image nginx skip auto-tambah "listen [::]:80" (cuma dilakukan kalau
# default.conf masih bawaan image), jadi nginx cuma listen IPv4. Kalau
# healthcheck resolve "localhost" ke ::1 duluan (perilaku musl libc di
# alpine), wget dapat connection refused walau nginx-nya sehat.
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1/ >/dev/null || exit 1
