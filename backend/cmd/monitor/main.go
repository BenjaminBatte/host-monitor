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

	"github.com/BenjaminBatte/host-monitor/internal/config"
	handlers "github.com/BenjaminBatte/host-monitor/internal/handlers"
	"github.com/BenjaminBatte/host-monitor/internal/services"
	"github.com/BenjaminBatte/host-monitor/internal/storage"
	ws "github.com/BenjaminBatte/host-monitor/pkg/websocket"
	"github.com/joho/godotenv"
)

type Config struct {
	Hosts    []string
	Port     int
	Interval time.Duration
	WSPort   string
}

func parseFlags() *Config {
	hosts := flag.String("hosts", "", "Comma-separated list of hosts to monitor")
	port := flag.Int("port", 80, "Port to connect to (simulates ping)")
	interval := flag.Duration("interval", 5*time.Second, "Interval between checks")
	wsPort := flag.String("ws-port", ":9090", "HTTP/WebSocket server port (e.g., :9090)")
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

func persistLoop(ctx context.Context, monitor *services.MonitorService, db *storage.DB) {
	events := monitor.Events()
	for {
		select {
		case <-ctx.Done():
			return
		case ev := <-events:
			var latencyPtr *int32
			if ev.LatencyMs >= 0 {
				l := int32(ev.LatencyMs)
				latencyPtr = &l
			}
			if err := db.InsertCheck(context.Background(), ev.Host, ev.Up, latencyPtr, float32(ev.PacketLoss), ev.CheckedAt); err != nil {
				fmt.Printf("INSERT failed host=%s up=%v latency=%v err=%v\n", ev.Host, ev.Up, ev.LatencyMs, err)
			}
		}
	}
}

// loads settings once and then every 10s
func startConfigReloader() {
	if err := config.LoadSettings(); err != nil {
		fmt.Printf("Failed to load settings: %v\n", err)
		os.Exit(1)
	}
	go func() {
		for {
			time.Sleep(10 * time.Second)
			if err := config.LoadSettings(); err != nil {
				fmt.Printf("Error reloading settings: %v\n", err)
			}
		}
	}()
}

func Run(ctx context.Context, cfg *Config) {
	_ = godotenv.Load()

	// Optional DB
	var db *storage.DB
	if os.Getenv("DB_URL") != "" || os.Getenv("DB_HOST") != "" {
		d, err := storage.New(context.Background())
		if err != nil {
			panic(err)
		}
		db = d
		defer db.Pool.Close()
		fmt.Println("[DB] Connected")
	} else {
		fmt.Println("[DB] Skipped (no DB_URL/DB_HOST set)")
	}

	// Monitor + WS server
	monitor := services.NewMonitorService(cfg.Hosts)
	server := ws.NewWebSocketServer(monitor.GetMetricsStore())

	mux := http.NewServeMux()
	mux.HandleFunc("/ws", server.HandleConnections)
	mux.HandleFunc("/api/threshold", handlers.ThresholdHandler)

	// Serve embedded Angular app at "/"
	spa, err := handlers.SPA()
	if err != nil {
		panic(err)
	}
	mux.Handle("/", spa)

	httpServer := &http.Server{
		Addr:    "0.0.0.0" + cfg.WSPort,
		Handler: mux,
	}

	// broadcast metrics to WS clients
	go server.StartBroadcasting(ctx)

	// start HTTP server
	go func() {
		fmt.Printf("HTTP server started on %s\n", cfg.WSPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// start monitoring loop
	fmt.Printf("Monitoring hosts: %v every %v on port %d\n", cfg.Hosts, cfg.Interval, cfg.Port)
	go monitor.Start(cfg.Port, cfg.Interval)

	// persist to DB if enabled
	if db != nil {
		go persistLoop(ctx, monitor, db)
	}

	<-ctx.Done()
	fmt.Println("\n[Shutdown] Signal received. Cleaning up...")
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
	startConfigReloader()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	Run(ctx, cfg)
}
