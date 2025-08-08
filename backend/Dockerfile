FROM golang:1.21-alpine

WORKDIR /app

# Install bash (for running scripts)
RUN apk add --no-cache bash

COPY . .

RUN go build -o host-monitor ./cmd/monitor

# Expose backend port (for WebSocket)
EXPOSE 8080

# Default to running the ping-many script
RUN chmod +x ./scripts/ping-many.sh
CMD ["./scripts/ping-many.sh"]
