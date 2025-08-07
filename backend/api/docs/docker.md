# Docker Guide: Host Monitor (Linux)

This guide explains how to build and run the Host Monitor backend on a Linux system using Docker.



## 🐧 Prerequisites

- Docker installed on your Linux system
- You are running a kernel that supports `--network=host`
- ICMP (ping) is allowed and not blocked by firewall rules



## 🛠️ 1. Build the Docker Image

From the project root (`host-monitor/`), run:

```bash
docker build -t host-monitor .
This uses the Dockerfile in the root directory to compile the Go backend into a minimal container.

🚀 2. Run the Container (with Host Networking)
Use the following command to start the monitor service:

bash
Copy
Edit
docker run --rm \
  --network=host \
  host-monitor \
  --hosts=8.8.8.8,1.1.1.1 \
  --port=80 \
  --interval=5s
Why --network=host?
Required for accurate ICMP ping and port reachability checks.

Binds container networking directly to the host stack.

Ensures consistent ping and port behavior as if running natively.



🔁 Example Use Case
Monitor a set of production endpoints from a Linux server:

bash
Copy
Edit
docker run --rm \
  --network=host \
  host-monitor \
  --hosts=192.168.1.1,8.8.8.8 \
  --port=443 \
  --interval=10s
🧹 Cleanup
To stop the container (if running detached):

bash
Copy
Edit
docker ps
docker stop <container_id>
To remove the image:


docker rmi host-monitor