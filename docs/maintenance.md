---

title: Operations & Maintenance

---

[⬅ Back to Home](./) | [Overview](overview.md)  | [Architecture](architecture.md) | [Setup](setup.md) | [Usage](usage.md) | [Technology Choices](tech.md) | [Future Improvements](future_improvements.md)


# Operations & Maintenance

## Routine Tasks

* **Change hosts / thresholds:** edit `config/settings.json`

  * `hosts`: e.g. `["8.8.8.8","1.1.1.1"]`
  * `latencyThresholdMs`: e.g. `70`
* The backend auto-reloads settings (within \~10s). If changes don’t appear, restart the service.

## Service Management (systemd)

```bash
# status & logs
sudo systemctl status host-monitor
sudo journalctl -u host-monitor -n 200 -f

# restart & enable on boot
sudo systemctl restart host-monitor
sudo systemctl enable host-monitor

# stop / disable
sudo systemctl stop host-monitor
sudo systemctl disable host-monitor
```

## Quick Diagnostics

```bash
# open port (Ubuntu)
sudo ufw allow 9090/tcp

# is the app listening?
ss -ltnp | grep :9090 || sudo lsof -iTCP:9090 -sTCP:LISTEN

# test locally
curl -v http://127.0.0.1:9090
```

## Common Issues

* **Permission denied running script** → `chmod +x backend/scripts/ping-many.sh`
* **Service points to /opt but repo elsewhere** → update paths:

```bash
APP=/home/ubuntu/host-monitor
sudo sed -i "s|/opt/host-monitor|$APP|g" /etc/systemd/system/host-monitor.service
sed -i "s|^cd /opt/host-monitor/backend$|cd $APP/backend|" "$APP/backend/scripts/ping-many.sh"
sudo systemctl daemon-reload && sudo systemctl restart host-monitor
```

## Optional: Docker

```bash
# view logs
docker logs -f host-monitor

# restart
docker restart host-monitor

# rebuild backend (Makefile)
make build
```

[Overview](overview.md) | [Architecture](architecture.md) | [Setup](setup.md) | [Usage](usage.md) | [Technology Choices](tech.md) | [Future Improvements](future_improvements.md)

---

<sub>© 2025 Host Monitor • <a href="https://github.com/BenjaminBatte/host-monitor">GitHub Repo</a></sub>
