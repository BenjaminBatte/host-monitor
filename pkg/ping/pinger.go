package ping

import (
	"fmt"
	"net"
	"time"
)

type PingResult struct {
	Host    string
	Latency time.Duration
	Error   error
}

func CheckHost(host string, port int) PingResult {
	address := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	latency := time.Since(start)

	if err != nil {
		return PingResult{Host: host, Latency: 0, Error: err}
	}
	conn.Close()
	return PingResult{Host: host, Latency: latency, Error: nil}
}
