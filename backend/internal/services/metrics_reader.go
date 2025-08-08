package services

import "github.com/BenjaminBatte/host-monitor/internal/models"

type MetricsReader interface {
	All() map[string]models.HostMetrics
}
