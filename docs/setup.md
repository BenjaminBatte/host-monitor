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
sudo apt install -y git golang-go 

```

---

## 2. Clone & Build (Makefile recommended)

```bash
git clone https://github.com/BenjaminBatte/host-monitor.git
cd host-monitor
cd backend
make build   # builds ./host-monitor from ./backend/cmd/monitor
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

---

## 4. Install the systemd Unit

Using the Makefile:

```bash
make install-service  # installs unit, daemon-reload, enables & starts
```

---

## 5. Configure Hosts/Flags via Drop-in (Recommended)

Instead of editing the unit file directly:

```bash
sudo systemctl edit host-monitor
```

---

## 6. Start & Enable Service

If you used `make install-service`, the service is already enabled and started.

**Otherwise:**

```bash
sudo systemctl enable --now host-monitor
```

---

## 7. Verify & View Logs

```bash
make status
# or
systemctl status host-monitor
journalctl -u host-monitor -e -f
```

---



## 8. UI (Dashboard) Options

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

---

[⬅ Back to Home](./) | [Next → Usage](usage.md)
