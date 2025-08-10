package storage

import (
	"context"
	"time"
)

// InsertCheck writes one row into public.checks.
// - host: domain name as TEXT (no FK)
// - up: current status
// - latency: nil when DOWN
// - loss: percentage 0..100
// - at: timestamp of the check
func (db *DB) InsertCheck(ctx context.Context, host string, up bool, latency *int32, loss float32, at time.Time) error {
	_, err := db.Pool.Exec(ctx, `
		INSERT INTO public.checks (host, up, latency_ms, packet_loss, checked_at)
		VALUES ($1, $2, $3, $4, $5)
	`, host, up, latency, loss, at)
	return err
}
