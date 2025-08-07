package services

import (
	"fmt"
	"net"
	"time"
)

type MonitorService struct {
	hosts   []string
	metrics *MetricsStore
}

func NewMonitorService(hosts []string) *MonitorService {
	return &MonitorService{
		hosts:   hosts,
		metrics: NewMetricsStore(),
	}
}

func (m *MonitorService) CheckHost(host string, port int) (string, time.Duration, error) {
	address := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	latency := time.Since(start)

	if err != nil {
		return host, 0, err
	}
	conn.Close()
	return host, latency, nil
}

func (m *MonitorService) Start(port int, interval time.Duration) {
	for _, host := range m.hosts {
		go func(h string) {
			for {
				address := fmt.Sprintf("%s:%d", h, port)
				start := time.Now()
				conn, err := net.DialTimeout("tcp", address, 2*time.Second)
				duration := time.Since(start)

				if err != nil {
					fmt.Printf("[%s] DOWN (error: %v)\n", h, err)
					m.metrics.Update(h, 0, false)
				} else {
					fmt.Printf("[%s] UP (latency: %v)\n", h, duration)
					m.metrics.Update(h, duration, true)
					conn.Close()
				}

				time.Sleep(interval)
			}
		}(host)
	}
	select {}
}

func (m *MonitorService) GetMetricsStore() *MetricsStore {
	return m.metrics
}
