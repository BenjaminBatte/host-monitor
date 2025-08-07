# Host Monitor System Architecture

## Overview

The Host Monitor system is designed to continuously monitor the availability and performance of specified hosts and display the results in real time on a web-based dashboard. It consists of:

- A Go-based backend for network monitoring and WebSocket communication
- An Angular frontend that receives live data and visualizes it
- Docker and systemd support for deployment and management



## 🧱 Components

### 1. **Monitor Service (Go)**
- Accepts `--hosts`, `--port`, and `--interval` as CLI flags
- Uses ICMP pings to determine:
  - Host status (up/down)
  - Latency (in ms)
  - Uptime percentage
  - Packet loss
- Tracks metrics using in-memory data structures
- Broadcasts metrics to all connected WebSocket clients every second

### 2. **WebSocket Server**
- Listens at `ws://<host>:8080/ws`
- Pushes updated metrics to clients in real time
- Handles concurrent connections using goroutines and channels

### 3. **Angular UI**
- Connects to the WebSocket endpoint
- Parses incoming JSON payloads
- Displays:
  - Host status (color-coded)
  - Latency, uptime %, packet loss
  - Visuals via tables, charts, or dashboards

---

## 🔄 Data Flow

```text
[Go CLI Ping Service] 
        ↓
[Metrics Aggregator]
        ↓
[WebSocket Server]
        ↓
[Angular UI Client]
⚙️ Deployment Options
Systemd: Use host-monitor.service to run as a background Linux service

Docker: Containerize the backend for portability

Kubernetes (Optional): Helm charts or manifest files included in deployments/kubernetes/

📁 Directory Structure (Key Folders)

host-monitor/
├── cmd/monitor/         # Go entrypoint
├── internal/            # Config, models, services
├── pkg/ping/            # Ping logic
├── pkg/websocket/       # WebSocket logic
├── ui/                  # Angular frontend
├── deployments/         # Docker/K8s files
├── docs/                # Documentation files


Design Decisions
Language: Go for performance, concurrency, and static compilation

WebSocket: Chosen for low-latency, push-based communication

Angular: Component-driven UI for real-time dashboards

Separation of Concerns: Clean package and module organization for testability and maintainability

Extensibility
This architecture allows for future enhancements such as:

Database persistence

Authentication for UI

REST API endpoints for historical data

Custom alerting (e.g., Slack, email)