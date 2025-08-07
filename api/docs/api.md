# WebSocket API Documentation

## Endpoint
ws://localhost:8080/ws

## Description
This WebSocket endpoint streams real-time host metrics to all connected clients. The data is updated and pushed every second, allowing the frontend to display up-to-date network statistics.

## Message Format

Each message is a JSON object where each key is an IP address or hostname, and the value is a set of metrics for that host.

### Example Payload
```json
{
  "1.1.1.1": {
    "up": true,
    "latency": 15,
    "uptime": 99.5,
    "packetLoss": 0
  },
  "8.8.8.8": {
    "up": false,
    "latency": null,
    "uptime": 72.3,
    "packetLoss": 12.5
  }
}
Field Descriptions
Field	Type	Description
up	boolean	true if the host is reachable, else false
latency	number or null	Round-trip time in milliseconds, null if host is down
uptime	number	Estimated uptime percentage over monitoring time
packetLoss	number	Packet loss percentage

Behavior
The backend broadcasts a full snapshot every second to all clients.

If a host goes down, latency is set to null.

Clients can connect at any time and will begin receiving the latest snapshot.