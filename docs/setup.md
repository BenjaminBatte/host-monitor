---

## title: Setup

---

[⬅ Back to Home](./) | [Next → Usage](usage.md)

# Setup

## 1. Prerequisites

* Ubuntu/Debian/CentOS with **systemd**
* **Go 1.22+** (to build) or prebuilt binary
* **PostgreSQL** (optional, for persistence)
* **Docker** (optional, for containerized deployment)

> Fresh Ubuntu/Debian packages:

```bash
sudo apt update
sudo apt install -y git golang-go postgresql postgresql-client ufw
# If not using Postgres, you can skip the two postgresql* packages
```

---

## 2. Clone & Build (Makefile recommended)

```bash
git clone https://github.com/BenjaminBatte/host-monitor.git
cd host-monitor
make build   # builds ./host-monitor from ./backend/cmd/monitor
```

**Alternative (manual):**

```bash
# Build backend binary
go build -o host-monitor ./backend/cmd/monitor
```

---

## 3. Install Binary & Create Service User

Create the service user (needed because the unit runs as `hostmonitor`):

```bash
sudo useradd --system --no-create-home --shell /usr/sbin/nologin hostmonitor
```

Install the binary (Makefile):

```bash
make install  # copies to /usr/local/bin/host-monitor
```

**Alternative (manual):**

```bash
sudo install -m 0755 host-monitor /usr/local/bin/host-monitor
```

---

## 4. (Optional) PostgreSQL Setup

If you want persistence, create a database and run the migration:

```bash
sudo -u postgres psql <<'SQL'
CREATE DATABASE hostmonitor;
CREATE USER hostmon WITH PASSWORD 'changeme';
GRANT ALL PRIVILEGES ON DATABASE hostmonitor TO hostmon;
SQL

export DB_URL='postgres://hostmon:changeme@localhost:5432/hostmonitor?sslmode=disable'
# Adjust path if your repo layout differs
psql "$DB_URL" -v ON_ERROR_STOP=1 -f backend/migrations/001_init.sql
```

---

## 5. (Optional) Environment File

Keep DB credentials and environment variables in a secure file:

```bash
sudo mkdir -p /etc/host-monitor
sudo tee /etc/host-monitor/host-monitor.env >/dev/null <<'EOF'
DB_URL=postgres://hostmon:changeme@localhost:5432/hostmonitor?sslmode=disable
GO_ENV=production
# Add other env vars here
EOF
sudo chmod 600 /etc/host-monitor/host-monitor.env
sudo chown root:root /etc/host-monitor/host-monitor.env
```

> In the unit file / drop-in, include:
> `EnvironmentFile=/etc/host-monitor/host-monitor.env`

---

## 6. Install the systemd Unit

Using the Makefile:

```bash
make install-service  # installs unit, daemon-reload, enables & starts
```

**Alternative (manual):**

```bash
sudo cp backend/deployments/host-monitor.service /etc/systemd/system/host-monitor.service
sudo systemctl daemon-reload
sudo systemctl enable --now host-monitor
```

The provided unit:

* Runs as `User=hostmonitor`, `Group=hostmonitor`
* Uses `WorkingDirectory=/usr/local/bin`
* Starts with `ExecStart=/usr/local/bin/host-monitor --hosts=1.1.1.1,8.8.8.8 --port=80 --interval=5s`
* Includes security hardening and journald logging

---

## 7. Configure Hosts/Flags via Drop-in (Recommended)

Instead of editing the unit file directly:

```bash
sudo systemctl edit host-monitor
```

Example drop-in override (recommended):

```ini
[Service]
EnvironmentFile=/etc/host-monitor/host-monitor.env
ExecStart=
ExecStart=/usr/local/bin/host-monitor \
  --hosts=8.8.8.8,1.1.1.1 \
  --port=80 \
  --interval=5s \
  --ws-port=:9090
```

> `systemctl edit` writes the drop-in and reloads systemd automatically. Restart to apply:

```bash
sudo systemctl restart host-monitor
```

---

## 8. Start & Enable Service

If you used `make install-service`, the service is already enabled and started.

**Otherwise:**

```bash
sudo systemctl enable --now host-monitor
```

---

## 9. Verify & View Logs

```bash
make status
# or
systemctl status host-monitor
journalctl -u host-monitor -e -f
```

---

## 10. Firewall (If Needed)

If WebSocket/API listens on `:9090` (as in the drop-in above):

```bash
sudo ufw allow 9090/tcp
```

If your API listens on another port (e.g., `:8180`), allow that instead:

```bash
sudo ufw allow 9090/tcp
```

---

## 11. ICMP vs TCP (Capabilities)

If the monitor uses **raw ICMP ping** in the future, add capability for a non-root service:

```bash
sudo setcap cap_net_raw+ep /usr/local/bin/host-monitor
```

Or in the unit/drop-in add:

```ini
AmbientCapabilities=CAP_NET_RAW
CapabilityBoundingSet=CAP_NET_RAW
```

> With the default TCP port check approach, this is **not required**.

---

## 12. UI (Dashboard) Options

This setup runs the **backend** service. To run the **Angular UI** on the same machine:

* **Dev**

  ```bash
  cd ui
  npm ci
  npm run start
  # Configure the UI to point to ws://<server-ip>:9090
  ```

* **Prod (serve built UI)**

  * Build: `npm run build`
  * Serve via Nginx/Apache or any static server (or use Docker compose in this repo).

---

## 13. Updating

**Makefile:**

```bash
make build
make install
sudo systemctl restart host-monitor
```

**Alternative (manual):**

```bash
# From repo root
go build -o host-monitor ./backend/cmd/monitor
sudo install -m 0755 host-monitor /usr/local/bin/host-monitor
sudo systemctl restart host-monitor
```

---

## 14. SELinux (CentOS/RHEL)

If SELinux is enforcing and you see denials, check logs and apply the suggested fixes:

```bash
sudo ausearch -m avc -ts recent
```

---

[⬅ Back to Home](./) | [Next → Usage](usage.md)
