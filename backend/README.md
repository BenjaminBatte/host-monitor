# Host Monitor – Backend

## Project Information

Full documentation:
[https://benjaminbatte.github.io/host-monitor/](https://benjaminbatte.github.io/host-monitor/)

## ⚡ Quick Start

**Option 1 – Run directly**

```bash
go run ./cmd/monitor \
  --hosts=8.8.8.8,1.1.1.1 \
  --port=80 \
  --interval=5s
```

**Option 2 – Use helper script**

```bash
chmod +x scripts/ping-many.sh
./scripts/ping-many.sh
```

**Flags:**

* `--hosts` — Comma-separated list of IPs or hostnames
* `--port` — TCP port to check for reachability
* `--interval` — Frequency of checks (e.g., 5s, 10s)

**WebSocket endpoint:**

```
ws://localhost:9090/ws
```

---

## 🐳 Docker (no Go toolchain required)

**Build**

```bash
docker build -t host-monitor-backend .
```

**Run (host network)**

```bash
docker run --rm \
  --network=host \
  -e DB_URL="postgres://postgres:2020@localhost:5432/postgres?sslmode=disable" \
  host-monitor-backend \
  --hosts=8.8.8.8,1.1.1.1 \
  --port=80 \
  --interval=5s
```

> If you’re on macOS/Windows with Docker Desktop, replace `localhost` in `DB_URL` with `host.docker.internal`.

**Run (mapped ports)**

```bash
docker run --rm -p 9090:9090 \
  -e DB_URL="postgres://postgres:2020@host.docker.internal:5432/postgres?sslmode=disable" \
  host-monitor-backend \
  --hosts=8.8.8.8,1.1.1.1 \
  --port=80 \
  --interval=5s
```

---

## ⚙️ Systemd (Linux service)

**Unit file** (`/etc/systemd/system/host-monitor.service`):

```ini
[Unit]
Description=Host Monitor Backend
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
Environment=DB_URL=postgres://postgres:2020@localhost:5432/postgres?sslmode=disable
WorkingDirectory=/opt/host-monitor/backend
ExecStart=/usr/local/bin/host-monitor \
  --hosts=8.8.8.8,1.1.1.1 \
  --port=80 \
  --interval=5s
Restart=always
RestartSec=3
User=hostmon
Group=hostmon

[Install]
WantedBy=multi-user.target
```

**Install & start**

```bash
sudo cp host-monitor /usr/local/bin/
sudo useradd -r -s /usr/sbin/nologin hostmon || true
sudo mkdir -p /opt/host-monitor/backend
sudo chown -R hostmon:hostmon /opt/host-monitor
sudo systemctl daemon-reload
sudo systemctl enable host-monitor
sudo systemctl start host-monitor
sudo systemctl status host-monitor --no-pager
```

