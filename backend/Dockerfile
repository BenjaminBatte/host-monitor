
FROM golang:1.24.6-alpine

# Set working directory
WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy all source code
COPY . .

# Build the Go binary
RUN go build -o host-monitor ./cmd/monitor

# Command to run the monitor
CMD ["./host-monitor", "--hosts=1.1.1.1,8.8.8.8", "--port=80", "--interval=5s"]
