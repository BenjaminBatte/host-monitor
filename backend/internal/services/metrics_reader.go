package services

import "github.com/BenjaminBatte/host-monitor/internal/models"

// MetricsReader is an interface for anything that can return host metrics.
type MetricsReader interface {
	All() map[string]models.HostMetrics
}
