---
## title: Setup
---

[⬅ Back to Home](./) | [Next → Usage](usage.md)

# Setup

## 1) Prerequisites

* Ubuntu/Debian/CentOS with **systemd**
* **Go 1.22+** (to build) or use a prebuilt binary

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

## 3) One-command install (recommended)

This builds the Go binary, copies the bundled UI to `/etc/host-monitor/web/`, installs the systemd unit, reloads systemd, and starts the service.

```bash
make deploy
```

> After this, the dashboard is available at:

```
http://<server-ip>:9090
```

---

## 4) Manual install (if you prefer separate steps)

```bash
# Build and install binary
make build
sudo install -m 0755 host-monitor /usr/local/bin/host-monitor

# Copy bundled Angular UI to the path the unit expects
sudo mkdir -p /etc/host-monitor/web
sudo cp -R web/* /etc/host-monitor/web/

# Install + enable + start the systemd unit
make install-service
```

---

## 5) Verify & Logs

```bash
systemctl status host-monitor
journalctl -u host-monitor -e -f
```

Open the dashboard:

```
http://<server-ip>:9090
```

* UI is served from the backend
* WebSocket endpoint is available at `/ws`

---

## 6) Update to latest

```bash
cd host-monitor/backend
make build
sudo cp host-monitor /usr/local/bin/host-monitor
sudo systemctl restart host-monitor
```

---

## Optional: Database

No database is required. To persist checks to PostgreSQL, set `DB_URL` in:

```
/etc/host-monitor/host-monitor.env
```

and restart the service.

---

[⬅ Back to Home](./) | [Next → Usage](usage.md)
