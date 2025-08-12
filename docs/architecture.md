---

## title: Architecture

---

[⬅ Back to Home](./) | [Next → Usage](usage.md) | [→ Setup](setup.md)

---

# System Architecture

This page describes how Host Monitor is structured, how data flows through the system, and the key runtime behaviors.

## High‑Level Components

![Class Diagram — Key Components & Data](./diagrams/class-diagram-key-components-data.svg)

**What’s shown:** Angular UI, WebSocket server, Go services (`MonitorService`, `MetricsStore`), persistence (PostgreSQL), and core data models.

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

## Notes on Persistence

* Each check is written to `public.checks (host, up, latency_ms, packet_loss, checked_at)`.
* DB connectivity is optional for local/dev; when configured, events are persisted asynchronously.

---

[⬅ Back to Setup](setup.md) | [Next → Maintenance](maintenance.md)
---

<sub>© 2025 Host Monitor • <a href="https://github.com/BenjaminBatte/host-monitor">GitHub Repo</a></sub>
