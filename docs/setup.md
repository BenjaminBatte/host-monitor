---
title: Setup
---

[⬅ Back to Home](./) | [Start → Overview](overview.md) | [Architecture](architecture.md) | [Usage](usage.md) | [Maintenance](maintenance.md) | [Technology Choices](tech.md) | [Future Improvements](future_improvements.md)


# Setup

## 1) Prerequisites

Linux with systemd (Ubuntu/Debian/CentOS)
Go 1.22+

```bash
sudo apt update
sudo apt install -y git golang-go
cd /opt
```

## 2) Clone

```bash
git clone https://github.com/BenjaminBatte/host-monitor.git
cd host-monitor/backend
```

> **Note:** The default service and script hardcode `/opt/host-monitor`. Cloning into `/opt` (as above) means **no extra steps**.
>
> **Not cloning into ****`/opt`****? Quick fix (example path ****`/home/ubuntu/host-monitor`****):**
>
> ```bash
> APP=/home/ubuntu/host-monitor
> # Update service to point to your path
> sudo sed -i "s|/opt/host-monitor|$APP|g" /etc/systemd/system/host-monitor.service
> # Update the script's cd line
> sed -i "s|^cd /opt/host-monitor/backend$|cd $APP/backend|" "$APP/backend/scripts/ping-many.sh"
> # Reload + restart
> sudo systemctl daemon-reload && sudo systemctl restart host-monitor
> ```

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

The SPA UI is served directly by the backend (embedded via go\:embed)
WebSocket endpoint: /ws

## 4) Verify & Logs

```bash
systemctl status host-monitor
journalctl -u host-monitor -e -f
```

---

**If the page doesn’t load:**

```bash
sudo ufw allow 9090/tcp
```

## 5) Other OS (macOS/Windows/WSL)

If you’re **not on Linux** (so you’re not using systemd), you can run it directly:

1. Clone anywhere (no need for `/opt`).
2. Edit `backend/scripts/ping-many.sh` and **remove** this line:

   ```
   cd /opt/host-monitor/backend
   ```
3. From the `backend` folder, run:

   ```bash
   chmod +x ./scripts/ping-many.sh
   ./scripts/ping-many.sh
   ```
4. Open the dashboard:

   ```
   http://localhost:9090
   ```

 [Start → Architecture](architecture.md) | [Usage](usage.md) | [Maintenance](maintenance.md) | [Technology Choices](tech.md) | [Future Improvements](future_improvements.md)


---

<sub>© 2025 Host Monitor • <a href="https://github.com/BenjaminBatte/host-monitor">GitHub Repo</a></sub>
