<<<<<<< HEAD
# Host Monitor

A real-time network monitoring tool built in **Go**, designed to track host availability, latency, uptime, and packet loss — with live updates via **WebSocket** to a dynamic **Angular** frontend.

---

## Features

 Monitor multiple hosts simultaneously
Track latency, uptime percentage, and packet loss
WebSocket server for real-time metrics streaming
Angular UI with live charts and status dashboard
Systemd service support
Docker-compatible for portable deployment
Built in **Go** 

---

## CLI Usage

Run the backend with:

```bash
go run ./cmd/monitor \
  --hosts=8.8.8.8,1.1.1.1 \
  --port=80 \
  --interval=5s
--hosts: Comma-separated list of IPs or hostnames

--port: TCP port to check for reachability

--interval: Frequency of pinging hosts (e.g., 5s, 10s)

🌐 WebSocket API
Endpoint: ws://localhost:8080/ws

Pushes JSON data every second with host metrics

See docs/api.md for full schema

 Example Output (UI)
json

{
  "8.8.8.8": {
    "up": true,
    "latency": 22,
    "uptime": 99.9,
    "packetLoss": 0
  }
}

Project Structure
bash
Copy
Edit
host-monitor/
├── cmd/monitor/         # CLI entrypoint
├── internal/            # Core logic (models, services)
├── pkg/ping/            # ICMP ping
├── pkg/websocket/       # WebSocket broadcasting
├── ui/                  # Angular dashboard
├── docs/                # Documentation
├── deployments/         # Docker + Kubernetes configs
├── Makefile             # Build/test commands
├── Dockerfile
├── host-monitor.service # systemd unit file
└── README.md

Docker (Linux Only)
docker build -t host-monitor .
docker run --rm --network=host host-monitor \
  --hosts=8.8.8.8,1.1.1.1 \
  --port=80 \
  --interval=5s
See docs/docker.md

Systemd Service (Linux)
Install the binary and .service file, then:
sudo systemctl daemon-reexec
sudo systemctl enable host-monitor
sudo systemctl start host-monitor
See docs/system.md
=======


Frontend UI

The Angular-based dashboard connects to the backend WebSocket server and visualizes host metrics in real time.

To run the frontend:

```bash
cd ui/
ng serve
Then open http://localhost:4200/ in your browser.

The UI expects the WebSocket server to be available at:
ws://localhost:8080/
Use ng build to compile for production.

Documentation
WebSocket API

Architecture Overview

Docker Guide

Systemd Setup
>>>>>>> 81b8f9f (Backend)



## ⚡ Quick Start (Manual)

1. Make the helper script executable:
```bash
chmod +x scripts/ping-many.sh

Optionally run the script (example usage):
./scripts/ping-many.sh

Or run the monitor directly with flags:
./host-monitor --hosts=1.1.1.1,8.8.8.8 --port=80 --interval=5s