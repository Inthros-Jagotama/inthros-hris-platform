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

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost/ >/dev/null || exit 1
