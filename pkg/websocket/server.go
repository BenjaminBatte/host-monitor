package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/BenjaminBatte/host-monitor/internal/models"
	"github.com/BenjaminBatte/host-monitor/internal/services"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for dev/demo
	},
}

type WebSocketServer struct {
	clients   map[*websocket.Conn]bool
	metrics   *services.MetricsStore
	broadcast chan []byte
	mu        sync.Mutex
}

func NewWebSocketServer(metrics *services.MetricsStore) *WebSocketServer {
	return &WebSocketServer{
		clients:   make(map[*websocket.Conn]bool),
		metrics:   metrics,
		broadcast: make(chan []byte),
	}
}

func (ws *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	ws.mu.Lock()
	ws.clients[conn] = true
	ws.mu.Unlock()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			ws.mu.Lock()
			delete(ws.clients, conn)
			ws.mu.Unlock()
			break
		}
	}
}

func (ws *WebSocketServer) StartBroadcasting() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		rawMetrics := ws.metrics.All()
		transformed := make(map[string]models.HostMetricsDTO)

		for host, m := range rawMetrics {
			transformed[host] = m.ToDTO()
		}

		jsonData, err := json.Marshal(transformed)
		if err != nil {
			log.Printf("Failed to marshal metrics: %v", err)
			continue
		}

		ws.mu.Lock()
		for client := range ws.clients {
			err := client.WriteMessage(websocket.TextMessage, jsonData)
			if err != nil {
				log.Printf("WebSocket send error: %v", err)
				client.Close()
				delete(ws.clients, client)
			}
		}
		ws.mu.Unlock()
	}
}
