                                                                    [⬅ Back to Home](./) | [Next → Usage](usage.md)
---

## title: Architecture

<div style="text-align: right;">
  <a href="overview.md">⬅ Back to Overview</a> | <a href="setup.md">Next → Setup</a>
</div>

# Architecture

## Backend

The backend is a high-performance Go service that:

* Concurrently pings multiple hosts at a configurable interval using goroutines.
* Tracks and aggregates key metrics — latency, uptime percentage, and packet loss.
* Uses WebSockets to push updates instantly to all connected clients, ensuring real-time monitoring without polling overhead.
* Reads configuration dynamically (hosts, latency threshold) so updates can be applied without restarts.

## Frontend

The Angular-based dashboard is optimized for clarity and responsiveness:

* Displays real-time host status cards with latency, uptime, and packet loss indicators.
* Includes interactive charts for visualizing uptime distribution and latency trends.
* Provides a settings panel for host management and latency threshold configuration.
* Offers CSV export for historical data and copy-to-clipboard for quick sharing.
* Subscribes to WebSocket streams for seamless live updates.

## Deployment

The system supports flexible deployment approaches:

* **Docker** for portable, reproducible builds across environments.
* **systemd** for persistent, auto-starting services on Linux servers.
* **Kubernetes** manifests for cloud-native scaling and orchestration.
* Can serve the Angular UI via Nginx or any static file server.

[⬅ Back to Home](./) | [Next → Usage](usage.md)