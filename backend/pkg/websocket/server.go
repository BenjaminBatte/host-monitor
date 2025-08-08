package websocket

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/BenjaminBatte/host-monitor/internal/models"
	"github.com/BenjaminBatte/host-monitor/internal/services"
	"github.com/gorilla/websocket"
)

type WebSocketServer struct {
	metricsReader services.MetricsReader
	clients       map[*websocket.Conn]bool
	mu            sync.Mutex
}

func NewWebSocketServer(reader services.MetricsReader) *WebSocketServer {
	return &WebSocketServer{
		metricsReader: reader,
		clients:       make(map[*websocket.Conn]bool),
	}
}

func (s *WebSocketServer) HandleConnections(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("WebSocket upgrade failed: %v\n", err)
		return
	}

	s.mu.Lock()
	s.clients[conn] = true
	s.mu.Unlock()
}

func (s *WebSocketServer) StartBroadcasting(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			raw := s.metricsReader.All()
			data := make(map[string]models.HostMetricsDTO)
			for k, v := range raw {
				data[k] = v.ToDTO()
			}

			s.mu.Lock()
			for client := range s.clients {
				err := client.WriteJSON(data)
				if err != nil {
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		}
	}
}
