package main

import (
	"context"
	"flag"
	"strings"
	"testing"
	"time"
)

func TestRunWithValidConfig(t *testing.T) {
	cfg := &Config{
		Hosts:    []string{"127.0.0.1"},
		Port:     80,
		Interval: 1 * time.Second,
		WSPort:   ":8090",
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go Run(ctx, cfg)

	time.Sleep(2 * time.Second)

	cancel()

	// Allow some time for the server to shut down gracefully
	time.Sleep(500 * time.Millisecond)
}

func TestParseFlags(t *testing.T) {
	args := []string{
		"-hosts=1.1.1.1,8.8.8.8",
		"-port=80",
		"-interval=2s",
		"-ws-port=:8085",
	}

	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	hosts := fs.String("hosts", "", "Comma-separated list of hosts to monitor")
	port := fs.Int("port", 80, "Port to connect to")
	interval := fs.Duration("interval", 5*time.Second, "Interval between checks")
	wsPort := fs.String("ws-port", ":8080", "WebSocket server port")

	if err := fs.Parse(args); err != nil {
		t.Fatalf("Failed to parse flags: %v", err)
	}

	cfg := &Config{
		Hosts:    strings.Split(*hosts, ","),
		Port:     *port,
		Interval: *interval,
		WSPort:   *wsPort,
	}

	if len(cfg.Hosts) != 2 {
		t.Fatalf("Expected 2 hosts, got %d", len(cfg.Hosts))
	}

	if cfg.Port != 80 || cfg.WSPort != ":8085" {
		t.Errorf("Unexpected values: port=%d, wsPort=%s", cfg.Port, cfg.WSPort)
	}

	if cfg.Interval != 2*time.Second {
		t.Errorf("Expected interval 2s, got %v", cfg.Interval)
	}
}
