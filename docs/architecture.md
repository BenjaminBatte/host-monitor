---
title: Architecture
---

[⬅ Back to Home](./) | [Next → Usage](usage.md) | [→ Setup](setup.md)

# Architecture

## Backend
The backend is a high-performance **Go** service designed for real-time network monitoring:

- **Concurrent host checks:** Uses goroutines to ping multiple hosts at a configurable interval.
- **Metric aggregation:** Tracks latency, uptime percentage, and packet loss.
- **Real-time updates:** Pushes data instantly to all connected clients via WebSockets — no polling required.
- **Dynamic configuration:** Supports live updates to hosts and latency threshold without restarting the service.

## Frontend
The **Angular**-based dashboard is optimized for clarity, speed, and responsiveness:

- **Live status cards:** Show latency, uptime, and packet loss for each monitored host.
- **Interactive charts:** Visualize uptime distribution and latency trends.
- **Settings panel:** Add/remove hosts and adjust latency threshold in real time.
- **Data export & sharing:** CSV download for historical data and one-click copy-to-clipboard for host details.
- **WebSocket integration:** Subscribes to backend streams for seamless live updates.

## Deployment
Flexible deployment options allow the system to run in multiple environments:

- **Docker** – Portable, reproducible builds for any platform.
- **systemd** – Persistent, auto-starting services on Linux.
- **Kubernetes** – Manifests for cloud-native scaling and orchestration.
- **Static hosting** – Serve the Angular UI via Nginx or any static file server.

[⬅ Back to Home](./) | [Next → Usage](usage.md)
