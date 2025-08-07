build:
	go build -o host-monitor ./cmd/monitor

install:
	sudo cp host-monitor /usr/local/bin/host-monitor

install-service:
	sudo cp deployments/host-monitor.service /etc/systemd/system/
	sudo systemctl daemon-reload
	sudo systemctl enable host-monitor
	sudo systemctl start host-monitor

status:
	sudo systemctl status host-monitor
