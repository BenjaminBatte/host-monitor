package services

import (
	"fmt"
	"net"
	"time"

	"github.com/BenjaminBatte/host-monitor/internal/config"
)

type MonitorService struct {
	hosts   []string
	metrics *MetricsStore
	events  chan CheckEvent
}

func NewMonitorService(hosts []string) *MonitorService {
	return &MonitorService{
		hosts:   hosts,
		metrics: NewMetricsStore(),
		events:  make(chan CheckEvent, 1024),
	}
}

func (m *MonitorService) Events() <-chan CheckEvent { return m.events }

func (m *MonitorService) CheckHost(host string, port int) (string, time.Duration, error) {
	address := fmt.Sprintf("[%s]:%d", host, port)

	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	latency := time.Since(start)

	if err != nil {
		return host, 0, err
	}
	if conn != nil {
		conn.Close()
	}
	return host, latency, nil
}

func (m *MonitorService) Start(port int, interval time.Duration) {
	for _, host := range m.hosts {
		go m.monitorHost(host, port, interval)
	}
	select {} // block forever
}

func (m *MonitorService) monitorHost(h string, port int, interval time.Duration) {
	for {
		address := fmt.Sprintf("[%s]:%d", h, port)
		start := time.Now()
		conn, err := net.DialTimeout("tcp", address, 2*time.Second)
		duration := time.Since(start)
		threshold := config.GetThreshold()

		up := (err == nil) && (duration.Milliseconds() <= int64(threshold))

		if up {
			fmt.Printf("[%s] UP (latency: %d ms, threshold: %d ms)\n", h, duration.Milliseconds(), threshold)
			m.metrics.Update(h, duration, true)
			if conn != nil {
				conn.Close()
			}
		} else {
			fmt.Printf("[%s] DOWN (latency: %d ms, threshold: %d ms, err: %v)\n", h, duration.Milliseconds(), threshold, err)
			m.metrics.Update(h, 0, false)
			if err == nil && conn != nil {
				conn.Close()
			}
		}

		// Compute current window packet-loss from the model after recording this probe.
		loss := 0.0
		if cur := m.metrics.Get(h); cur != nil {
			loss = cur.PacketLossWindow()
		}

		ev := CheckEvent{
			Host: h,
			Up:   up,
			LatencyMs: func() int {
				if up {
					return int(duration.Milliseconds())
				}
				return -1
			}(),
			PacketLoss: loss, // window loss %
			CheckedAt:  time.Now().UTC(),
		}
		select {
		case m.events <- ev:
		default:
			// drop event if channel is full to avoid blocking
		}

		time.Sleep(interval)
	}
}

// Expose metrics for HTTP handlers, etc.
func (m *MonitorService) GetMetricsStore() MetricsReader {
	return m.metrics
}
