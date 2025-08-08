
FROM golang:1.24.6-alpine


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o host-monitor ./cmd/monitor


CMD ["./host-monitor", "--hosts=1.1.1.1,8.8.8.8", "--port=80", "--interval=5s"]
