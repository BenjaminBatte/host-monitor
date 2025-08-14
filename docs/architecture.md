---

## title: Architecture

---

[⬅ Back to Home](./) | [Start → Overview](overview.md) | [Setup](setup.md) | [Usage](usage.md) | [Maintenance](maintenance.md) | [Technology Choices](tech.md) | [Future Improvements](future_improvements.md)


---

# System Architecture

This page describes how Host Monitor is structured, how data flows through the system, and the key runtime behaviors.

## High‑Level Components

![Class Diagram — Key Components & Data](./diagrams/class-diagram-key-components-data.svg)

**What’s shown:** Angular UI, WebSocket server, Go services (`MonitorService`, `MetricsStore`), persistence (PostgreSQL), and core data models. The Angular UI is **embedded into the Go binary** using `go:embed` and served by the backend (no separate `backend/web` directory).

## Detailed Domain & Services

![Class Diagram — Domain & Services (Detailed)](./diagrams/class-diagram-domain-services-detailed.svg)

**What’s shown:** Method‑level view of domain models (`HostMetrics`, DTO), services (`MonitorService`, `MetricsStore`), settings/threshold API, and relationships.

## Runtime Flows

### Real‑Time Host Monitoring (Host 1 … Host N)

![Sequence Diagram — Real‑Time Host Monitoring – Host](./diagrams/sequence-diagram-real-time-host-monitoring-host.svg)

**Flow:** Browser → Angular Dashboard → WebSocket → Backend → Pingers → DB → WebSocket → UI.

### Metrics Push Loop

![Sequence Diagram — Real‑Time Metrics Push](./diagrams/sequence-diagram-real-time-metrics-push.svg)

**Flow:** Pingers report → `MetricsService` aggregates → DB upsert → broadcast JSON over `/ws` → UI renders.

### Update Latency Threshold

![Sequence Diagram — Update Latency Threshold](./diagrams/sequence-diagram-update-latency-threshold.svg)

**Flow:** UI slider (debounced) → send `{ setThreshold }` over `/ws` → backend acks and applies → UI toast.

## Connection Behavior

### WebSocket Reconnect & Backoff

![Flowchart — WebSocket Reconnect & Backoff](./diagrams/flowchart-websocket-reconnect-backoff.svg)

**Strategy:** Exponential backoff (1s, 2s, 4s … max 30s). On reconnect, UI resubscribes and resumes live updates.

### Host Status State Machine

![State Diagram — Host Status](./diagrams/state-diagram-host-status.svg)

**Rules:**

* First success → **Up**
* Fail or latency > threshold → **Down**
* Success with latency ≤ threshold → **Up**

## Embedded UI (no `backend/web`)

The production Angular build is embedded in the Go binary. At build time we copy the UI artifacts into the backend and embed them via `go:embed`:

```text
ui/dist/host-monitor-ui/browser/  →  backend/internal/handlers/dist/browser/
```

Example (in `internal/handlers/static.go`):

```go
//go:embed dist/browser/*
var uiFS embed.FS
var staticFS, _ = fs.Sub(uiFS, "dist/browser")
// mux.Handle("/", http.FileServer(http.FS(staticFS)))
```

**Why:**

* Single self‑contained binary (simpler deploys: systemd/Docker/K8s).
* No drift between filesystem and served assets.
* Safer defaults (read‑only embedded FS).

## Project Structure

The Host Monitor repository is organized into the following main components:

```
host-monitor
.
├── backend                 # Go backend service
│   ├── cmd                 # Entry points (CLI / main)
│   ├── deployments         # systemd, Docker, Kubernetes manifests
│   ├── internal            # Core app modules (config, services, handlers, models, utils)
│   │   ├── handlers
│   │   │   └── dist
│   │   │       └── browser # Embedded Angular UI (copied from ui/dist/host-monitor-ui/browser)
│   ├── migrations          # SQL migrations for PostgreSQL
│   ├── pkg                 # Reusable packages (ping, websocket, etc.)
│   ├── scripts             # Helper scripts (e.g., ping-many.sh)
│   ├── Dockerfile          # Backend container build
│   ├── go.mod / go.sum     # Go dependencies
│   ├── host-monitor        # Compiled binary (after build; git-ignored)
│   ├── Makefile            # Build/install tasks
│   └── README.md
├── docs                    # Documentation site
│   ├── diagrams            # Architecture, sequence, and flow diagrams
│   ├── images              # Screenshots and visuals for docs
│   ├── architecture.md     # This page
│   ├── future_improvements.md
│   ├── index.md
│   ├── maintenance.md
│   ├── overview.md
│   ├── setup.md
│   ├── tech.md
│   └── usage.md
├── ui                      # Angular frontend
│   ├── dist                # Production build output
│   ├── node_modules        # Frontend dependencies
│   ├── public              # Static assets
│   ├── src                 # Angular components, services, styles
│   ├── angular.json        # Angular CLI config
│   ├── Dockerfile          # UI container build
│   ├── package*.json       # NPM dependencies
│   ├── proxy.conf.json     # Dev proxy settings
│   ├── README.md
│   └── tsconfig*.json      # TypeScript configs
├── docker-compose.dev.yml  # Multi-container dev environment
└── README.md               # Root project overview
```

**Key conventions:**

* **backend/** contains the Go services, WebSocket server, metrics logic, and optional DB integration.
* **ui/** is the Angular SPA, served separately in dev and **embedded into the backend in production** at `backend/internal/handlers/dist/browser`.
* **docs/** is a MkDocs-style documentation site (the `.md` files you’re reading now).
* **scripts/** has helper and deployment scripts (kept untracked in some cases for local overrides).
* **deployments/** contains systemd, Docker Compose, and Kubernetes specs for easy deployment.

## Notes on Persistence

* Each check is written to `public.checks (host, up, latency_ms, packet_loss, checked_at)`.
* DB connectivity is optional for local/dev; when configured, events are persisted asynchronously.

---

[Start → Overview](overview.md) | [Setup](setup.md) | [Usage](usage.md) | [Maintenance](maintenance.md) | [Technology Choices](tech.md) | [Future Improvements](future_improvements.md)

---

<sub>© 2025 Host Monitor • <a href="https://github.com/BenjaminBatte/host-monitor">GitHub Repo</a></sub>
