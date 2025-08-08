package models

import (
	"sync"
	"time"
)

type HostStatus string

const (
	StatusUp   HostStatus = "UP"
	StatusDown HostStatus = "DOWN"
)

type HostMetrics struct {
	Host           string
	LatencyHistory []time.Duration
	UpCount        int
	DownCount      int
	LastStatus     HostStatus
	LastChecked    time.Time
	mu             sync.Mutex
}

func (hm *HostMetrics) AddResult(latency time.Duration, status HostStatus) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	hm.LatencyHistory = append(hm.LatencyHistory, latency)
	if len(hm.LatencyHistory) > 100 {
		hm.LatencyHistory = hm.LatencyHistory[1:] // Keep only the last 100 entries
	}

	if status == StatusUp {
		hm.UpCount++
	} else {
		hm.DownCount++
	}

	hm.LastStatus = status
	hm.LastChecked = time.Now().UTC()
}

type HostMetricsDTO struct {
	Latency        float64   `json:"latency"`
	Up             bool      `json:"up"`
	TotalChecks    int       `json:"totalChecks"`
	SuccessCount   int       `json:"successCount"`
	LastChecked    string    `json:"lastChecked,omitempty"`
	LatencyHistory []float64 `json:"latencyHistory"`
}

func (hm *HostMetrics) ToDTO() HostMetricsDTO {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	latency := float64(0)
	if len(hm.LatencyHistory) > 0 {
		latency = float64(hm.LatencyHistory[len(hm.LatencyHistory)-1].Milliseconds())
	}

	latencyHistory := make([]float64, len(hm.LatencyHistory))
	for i, d := range hm.LatencyHistory {
		latencyHistory[i] = float64(d.Milliseconds())
	}

	total := hm.UpCount + hm.DownCount
	lastChecked := ""
	if !hm.LastChecked.IsZero() {
		lastChecked = hm.LastChecked.Format(time.RFC3339)
	}

	return HostMetricsDTO{
		Latency:        latency,
		Up:             hm.LastStatus == StatusUp,
		TotalChecks:    total,
		SuccessCount:   hm.UpCount,
		LastChecked:    lastChecked,
		LatencyHistory: latencyHistory,
	}
}
