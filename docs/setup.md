---
title: Setup
---

[⬅ Back to Home](./) | [Next → Usage](usage.md)

# Setup

## 1) Prerequisites

* Linux with **systemd** (Ubuntu/Debian/CentOS)
* **Go 1.22+**

```bash
sudo apt update
sudo apt install -y git golang-go
```

---

## 2) Clone

```bash
git clone https://github.com/BenjaminBatte/host-monitor.git
cd host-monitor/backend
```

---

## 3) Install & Run (embedded UI — no Angular needed)

```bash
# Build
make build

# Install binary
sudo install -m 0755 host-monitor /usr/local/bin/host-monitor

# Install + enable + start systemd unit
make install-service
```

Open:

```
http://<server-ip>:9090
```

* The SPA UI is served directly by the backend (embedded via `go:embed`)
* WebSocket endpoint: `/ws`

---

## 4) Verify & Logs

```bash
systemctl status host-monitor
journalctl -u host-monitor -e -f
```

---

[⬅ Back to Home](./) | [Next → Usage](usage.md)
