---

## title: Setup

[⬅ Back to Home](./) | [Next → Usage](usage.md)

# Setup

## Prerequisites

* Go 1.22+
* Node.js 18+ and Angular CLI
* Docker (optional, for containers)

---

## Local Development

### Backend (Go)

```bash
go run ./cmd/monitor --hosts=8.8.8.8,1.1.1.1 --port=80 --interval=5s
# Or from backend folder
./scripts/ping-many.sh
```

### Frontend (Angular)

```bash
cd ui
npm install
ng serve
```

**Proxy Config (`proxy.conf.json`)**

```json
{
  "/api": { "target": "http://localhost:8080", "secure": false, "changeOrigin": true, "logLevel": "debug" },
  "/ws":  { "target": "ws://localhost:8080",    "ws": true,    "secure": false, "changeOrigin": true, "logLevel": "debug" }
}
```

In `MetricsService`:

```ts
private readonly WS_URL = '/ws';
```

Run with proxy:

```bash
ng serve --proxy-config proxy.conf.json
```

---

## Docker

```bash
docker compose up --build
```

### Backend

```bash
docker build -t host-monitor-backend -f backend/Dockerfile .
```

### UI (if using Nginx-based UI Dockerfile)

```bash
docker build -t host-monitor-ui ./ui
docker run -d --name host-monitor-ui -p 80:80 host-monitor-ui
```

---

## Production Build (UI)

```bash
cd ui
ng build --configuration production
# Serve dist/ with nginx or any static server
```

---

## Notes

* Use normal hyphens `-` in flags (not en dashes `–`).
* For Docker v2, prefer `docker compose` (space), not `docker-compose`.

---

[⬅ Back to Home](./) | [Next → Usage](usage.md)
