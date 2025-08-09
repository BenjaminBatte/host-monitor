---
title: Architecture
---

# Architecture
**Backend:** Go service that pings hosts, tracks metrics (latency, uptime, packet loss), and sends data over WebSockets.  
**Frontend:** Angular dashboard with live charts, CSV export, and latency threshold config.  
**Deployment:** Docker, systemd, and optional Kubernetes manifests.
