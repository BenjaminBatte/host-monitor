# Systemd Setup: Host Monitor Service

This guide explains how to install and manage the Host Monitor as a background Linux service using `systemd`.


Prerequisites
- Linux system with `systemd` (e.g., Ubuntu, Debian, CentOS, Arch)
- Go binary or Docker image built and ready
- Root or sudo access

 1. Create the Systemd Unit File

Create the service file:

```bash
sudo nano /etc/systemd/system/host-monitor.service

Paste the following:
[Unit]
Description=Host Monitor Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/host-monitor \
    --hosts=8.8.8.8,1.1.1.1 \
    --port=80 \
    --interval=5s
Restart=on-failure
RestartSec=5
User=root
WorkingDirectory=/usr/local/bin

[Install]
WantedBy=multi-user.target
Adjust paths and flags as needed. Ensure the host-monitor binary is copied to /usr/local/bin/.

2. Enable and Start the Service
Reload systemd and start the service:
sudo systemctl daemon-reexec
sudo systemctl daemon-reload
sudo systemctl enable host-monitor
sudo systemctl start host-monitor

3. Check Service Status
sudo systemctl status host-monitor

4. View Logs
journalctl -u host-monitor -f
Use -f to follow logs in real-time.

 5. Restart or Stop the Service
Restart:
sudo systemctl restart host-monitor

Stop:
sudo systemctl stop host-monitor

Disable:
sudo systemctl disable host-monitor

Optional: Modify Without Rebuilding
If you want to change command-line flags without rebuilding:

Edit the service file:
sudo nano /etc/systemd/system/host-monitor.service

Reload and restart:
sudo systemctl daemon-reload
sudo systemctl restart host-monitor