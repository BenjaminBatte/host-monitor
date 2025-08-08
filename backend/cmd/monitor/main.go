package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/BenjaminBatte/host-monitor/internal/services"
	ws "github.com/BenjaminBatte/host-monitor/pkg/websocket"
)

type Config struct {
	Hosts    []string
	Port     int
	Interval time.Duration
	WSPort   string
}

// parseFlags parses command-line arguments and returns a Config object
func parseFlags() *Config {
	hosts := flag.String("hosts", "", "Comma-separated list of hosts to monitor")
	port := flag.Int("port", 80, "Port to connect to (simulates ping)")
	interval := flag.Duration("interval", 5*time.Second, "Interval between checks")
	wsPort := flag.String("ws-port", ":8080", "WebSocket server port")
	flag.Parse()

	if *hosts == "" {
		fmt.Println("Error: No hosts provided. Use -hosts to specify hosts.")
		return nil
	}

	return &Config{
		Hosts:    strings.Split(*hosts, ","),
		Port:     *port,
		Interval: *interval,
		WSPort:   *wsPort,
	}
}

// Run initializes and starts the monitor and WebSocket server
func Run(ctx context.Context, cfg *Config) {
	monitor := services.NewMonitorService(cfg.Hosts)
	server := ws.NewWebSocketServer(monitor.GetMetricsStore())

	// Create a mux and register the WebSocket handler
	mux := http.NewServeMux()
	mux.HandleFunc("/ws", server.HandleConnections)

	// Define an HTTP server with shutdown capability
	httpServer := &http.Server{
		Addr:    "0.0.0.0" + cfg.WSPort,
		Handler: mux,
	}

	// Start broadcasting (respects context)
	go server.StartBroadcasting(ctx)

	// Start the HTTP server
	go func() {
		fmt.Printf("WebSocket server started on %s\n", cfg.WSPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("WebSocket server error: %v\n", err)
		}
	}()

	// Start monitoring
	fmt.Printf("Monitoring hosts: %v every %v on port %d\n", cfg.Hosts, cfg.Interval, cfg.Port)
	go monitor.Start(cfg.Port, cfg.Interval)

	// Wait for Ctrl+C (context cancellation)
	<-ctx.Done()
	fmt.Println("\n[Shutdown] Signal received. Cleaning up...")

	// Shutdown the HTTP server gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("HTTP server shutdown error: %v\n", err)
	}

	fmt.Println("[Shutdown] Complete.")
}

func main() {
	cfg := parseFlags()
	if cfg == nil {
		return
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	Run(ctx, cfg)
}
