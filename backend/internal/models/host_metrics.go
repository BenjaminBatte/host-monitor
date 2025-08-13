package models

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

type HostStatus int

const (
	StatusUp HostStatus = iota
	StatusDown
)

func (s HostStatus) String() string {
	switch s {
	case StatusUp:
		return "UP"
	case StatusDown:
		return "DOWN"
	default:
		return "UNKNOWN"
	}
}

func (s HostStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

func (s *HostStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	switch str {
	case "UP":
		*s = StatusUp
	case "DOWN":
		*s = StatusDown
	default:
		return fmt.Errorf("invalid HostStatus: %s", str)
	}
	return nil
}

const windowSize = 100 // sliding window size for recent checks

type HostMetrics struct {
	Host           string
	LatencyHistory []time.Duration // last N latencies (ms)
	StatusHistory  []bool          // last N statuses: true=UP, false=DOWN

	UpCount     int
	DownCount   int
	LastStatus  HostStatus
	LastChecked time.Time

	mu sync.Mutex
}

func (hm *HostMetrics) AddResult(latency time.Duration, status HostStatus) {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	// Sliding windows
	hm.LatencyHistory = append(hm.LatencyHistory, latency)
	if len(hm.LatencyHistory) > windowSize {
		hm.LatencyHistory = hm.LatencyHistory[1:]
	}

	hm.StatusHistory = append(hm.StatusHistory, status == StatusUp)
	if len(hm.StatusHistory) > windowSize {
		hm.StatusHistory = hm.StatusHistory[1:]
	}

	// Lifetime counters
	if status == StatusUp {
		hm.UpCount++
	} else {
		hm.DownCount++
	}

	hm.LastStatus = status
	hm.LastChecked = time.Now().UTC()
}

// PacketLossWindow returns loss % over the recent window (StatusHistory).
func (hm *HostMetrics) PacketLossWindow() float64 {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	n := len(hm.StatusHistory)
	if n == 0 {
		return 0
	}
	fail := 0
	for _, ok := range hm.StatusHistory {
		if !ok {
			fail++
		}
	}
	return float64(fail) / float64(n) * 100.0
}

// PacketLossLifetime returns loss % since process start.
func (hm *HostMetrics) PacketLossLifetime() float64 {
	hm.mu.Lock()
	defer hm.mu.Unlock()

	total := hm.UpCount + hm.DownCount
	if total == 0 {
		return 0
	}
	return float64(hm.DownCount) / float64(total) * 100.0
}

type HostMetricsDTO struct {
	Latency           float64   `json:"latency"` // ms
	Up                bool      `json:"up"`
	TotalChecks       int       `json:"totalChecks"`
	SuccessCount      int       `json:"successCount"`
	LastChecked       string    `json:"lastChecked,omitempty"`
	LatencyHistory    []float64 `json:"latencyHistory"` // ms
	PacketLossWindow  float64   `json:"packetLossWindow"`
	PacketLossOverall float64   `json:"packetLossOverall"`
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

	fail := 0
	for _, ok := range hm.StatusHistory {
		if !ok {
			fail++
		}
	}
	winLoss := 0.0
	if n := len(hm.StatusHistory); n > 0 {
		winLoss = float64(fail) / float64(n) * 100.0
	}
	overallLoss := 0.0
	if total > 0 {
		overallLoss = float64(hm.DownCount) / float64(total) * 100.0
	}

	return HostMetricsDTO{
		Latency:           latency,
		Up:                hm.LastStatus == StatusUp,
		TotalChecks:       total,
		SuccessCount:      hm.UpCount,
		LastChecked:       lastChecked,
		LatencyHistory:    latencyHistory,
		PacketLossWindow:  winLoss,
		PacketLossOverall: overallLoss,
	}
}
