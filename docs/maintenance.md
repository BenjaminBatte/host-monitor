---

## title: Maintenance

[⬅ Back to Home](./) | [⬅ Back to Usage](usage.md)

# Maintenance

## Routine Tasks

* **Add/Remove Hosts**

  * Update `config/settings.json` (e.g., `"hosts": ["8.8.8.8","1.1.1.1"]`).
  * The backend reloads settings periodically (no restart needed). If changes don’t apply, restart the service.

* **Latency Threshold**

  * Modify `"latencyThresholdMs"` in `config/settings.json` or adjust via the UI.
  * The backend reloads this value automatically (used to mark hosts DOWN when latency exceeds the threshold).

---

## Service Management

### systemd (Linux) & Docker

```bash
# status / logs
sudo systemctl status host-monitor
journalctl -u host-monitor -n 200 -f

# restart / enable on boot
sudo systemctl restart host-monitor
sudo systemctl enable host-monitor

# view logs (Docker)
docker logs -f host-monitor

# restart container (Docker)
docker restart host-monitor

# pull latest image & restart (Docker)
docker pull <ACCOUNT_ID>.dkr.ecr.<REGION>.amazonaws.com/host-monitor:latest
docker stop host-monitor && docker rm host-monitor
docker run -d --name host-monitor -p 80:8080 <...> # with your flags

# or rebuild locally
make build        # builds backend
make ui-build     # builds Angular UI

# validate settings.json syntax
jq . config/settings.json >/dev/null
```

---

[⬅ Back to Home](./) | [⬅ Back to Usage](usage.md)
