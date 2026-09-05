# ===================================================================
# Stage 1: Build bundle Vite
# ===================================================================
FROM node:22-alpine AS builder

WORKDIR /app

# Copy manifest lebih dulu agar layer npm ci ter-cache selama dependency
# tidak berubah.
COPY package.json package-lock.json ./
RUN npm ci

COPY . .

# .env.development sengaja tidak ikut (di-ignore lewat .dockerignore):
# VITE_COMPANY_ID di sana hanya alat bantu login dev.
RUN npm run build

# ===================================================================
# Stage 2: Serve static via nginx
# ===================================================================
FROM nginx:1.27-alpine

# src/services/api.js memakai `baseURL: ''` — seluruh panggilan API berupa
# path relatif (/api/...), jadi bundle ini WAJIB diserve dari origin yang
# sama dengan backend. nginx.conf yang mem-proxy /api dan /uploads ke
# service `api` adalah yang membuat asumsi itu terpenuhi.
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /app/dist /usr/share/nginx/html

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget -qO- http://localhost/ >/dev/null || exit 1
