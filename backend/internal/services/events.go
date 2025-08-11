package services

import "time"

// Emitted after each host check.
type CheckEvent struct {
	Host       string
	Up         bool
	LatencyMs  int     // -1 when unknown/down
	PacketLoss float64 //to be used for future checks
	CheckedAt  time.Time
}
