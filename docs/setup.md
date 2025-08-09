---
title: Setup
---

# Setup

## Prerequisites
- Go 1.22+
- Node.js 18+ and Angular CLI
- Docker (optional, for containers)

---

## Local Development

### Backend (Go)
```bash
go run ./cmd/monitor --hosts=8.8.8.8,1.1.1.1 --port=80 --interval=5s
cd ui
npm install
ng serve
{
  "/api": { "target": "http://localhost:8080", "secure": false, "changeOrigin": true, "logLevel": "debug" },
  "/ws":  { "target": "ws://localhost:8080",    "ws": true,    "secure": false, "changeOrigin": true, "logLevel": "debug" }
}
// MetricsService
private readonly WS_URL = '/ws';
cd ui
ng serve --proxy-config proxy.conf.json
docker compose up --build
# Backend
docker build -t host-monitor-backend -f backend/Dockerfile .

# UI (if you have an Nginx-based UI Dockerfile)
docker build -t host-monitor-ui ./ui
docker run -d --name host-monitor-ui -p 80:80 host-monitor-ui
./scripts/ping-many.sh
cd ui
ng build --configuration production
# Serve dist/ with nginx or any static server

**Notes**
- Use normal hyphens `-` in flags (not en dashes `–`).
- For Docker v2, prefer `docker compose` (space), not `docker-compose`.

Want me to add a **“Next → Usage”** button at the bottom like we did for the Home button?
::contentReference[oaicite:0]{index=0}
```
[⬅ Back to Home](./)