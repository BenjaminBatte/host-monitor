package services

import (
	"sync"
	"time"

	"github.com/BenjaminBatte/host-monitor/internal/models"
)

// MetricsStore holds per-host metrics with concurrency safety.
type MetricsStore struct {
	mu      sync.RWMutex
	metrics map[string]*models.HostMetrics
}

// NewMetricsStore creates a new empty store.
func NewMetricsStore() *MetricsStore {
	return &MetricsStore{
		metrics: make(map[string]*models.HostMetrics),
	}
}

// Update records a new check result for a host.
func (ms *MetricsStore) Update(host string, latency time.Duration, isUp bool) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	m, exists := ms.metrics[host]
	if !exists {
		m = &models.HostMetrics{Host: host}
		ms.metrics[host] = m
	}

	status := models.StatusDown
	if isUp {
		status = models.StatusUp
	}

	m.AddResult(latency, status)
}

func (ms *MetricsStore) Get(host string) *models.HostMetrics {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	return ms.metrics[host]
}

func (ms *MetricsStore) All() map[string]*models.HostMetrics {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	result := make(map[string]*models.HostMetrics, len(ms.metrics))
	for k, v := range ms.metrics {
		result[k] = v
	}
	return result
}

func (ms *MetricsStore) SnapshotDTO() map[string]models.HostMetricsDTO {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	out := make(map[string]models.HostMetricsDTO, len(ms.metrics))
	for k, v := range ms.metrics {
		out[k] = v.ToDTO()
	}
	return out
}
