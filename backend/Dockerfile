# backend/Dockerfile
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o host-monitor ./cmd/monitor

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/host-monitor .
EXPOSE 9090
CMD ["./host-monitor", "--hosts=8.8.8.8,1.1.1.1", "--port=80", "--interval=5s", "--ws-port=:9090"]
