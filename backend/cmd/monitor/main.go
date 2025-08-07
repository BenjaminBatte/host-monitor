package main

import (
	"flag"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/BenjaminBatte/host-monitor/internal/services"
	ws "github.com/BenjaminBatte/host-monitor/pkg/websocket"
)

func main() {
	hosts := flag.String("hosts", "", "Comma-separated list of hosts to monitor")
	port := flag.Int("port", 80, "Port to connect to (simulates ping)")
	interval := flag.Duration("interval", 5*time.Second, "Interval between checks")
	wsPort := flag.String("ws-port", ":8080", "WebSocket server port")
	flag.Parse()

	if *hosts == "" {
		fmt.Println("Error: No hosts provided. Use -hosts to specify hosts.")
		return
	}

	hostList := strings.Split(*hosts, ",")
	monitor := services.NewMonitorService(hostList)

	// Start WebSocket server
	server := ws.NewWebSocketServer(monitor.GetMetricsStore())
	go func() {
		http.HandleFunc("/ws", server.HandleConnections)
		go server.StartBroadcasting()
		fmt.Printf("WebSocket server started on %s\n", *wsPort)
		err := http.ListenAndServe("0.0.0.0"+*wsPort, nil)
		if err != nil {
			panic(err)
		}
	}()

	monitor.Start(*port, *interval)
}
