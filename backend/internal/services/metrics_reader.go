package services

import "github.com/BenjaminBatte/host-monitor/internal/models"

type MetricsReader interface {
	Get(host string) *models.HostMetrics
	All() map[string]*models.HostMetrics
}
