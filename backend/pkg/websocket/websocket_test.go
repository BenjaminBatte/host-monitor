package websocket

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/BenjaminBatte/host-monitor/internal/models"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

// --- Mock MetricsReader --- //
type mockMetricsStore struct {
	mu      sync.Mutex
	metrics map[string]*models.HostMetrics
}

func (m *mockMetricsStore) All() map[string]models.HostMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	result := make(map[string]models.HostMetrics)
	for k, v := range m.metrics {
		result[k] = *v
	}
	return result
}

func newMockMetricsStore() *mockMetricsStore {
	return &mockMetricsStore{
		metrics: map[string]*models.HostMetrics{
			"127.0.0.1": {
				Host:           "127.0.0.1",
				LatencyHistory: []time.Duration{10 * time.Millisecond},
				UpCount:        1,
				DownCount:      0,
				LastStatus:     models.StatusUp,
				LastChecked:    time.Now(),
			},
		},
	}
}

// --- Test --- //
func TestWebSocketBroadcastsMetrics(t *testing.T) {
	store := newMockMetricsStore()
	server := NewWebSocketServer(store)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpServer := httptest.NewServer(http.HandlerFunc(server.HandleConnections))
	defer httpServer.Close()

	wsURL := "ws" + httpServer.URL[len("http"):]

	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer wsConn.Close()

	// Set timeout and wait for broadcast
	wsConn.SetReadDeadline(time.Now().Add(3 * time.Second))
	go server.StartBroadcasting(ctx)
	time.Sleep(1500 * time.Millisecond)

	_, msg, err := wsConn.ReadMessage()
	assert.NoError(t, err)
	assert.NotEmpty(t, msg, "Expected JSON metrics message")

	var result map[string]models.HostMetricsDTO
	err = json.Unmarshal(msg, &result)
	assert.NoError(t, err, "Should parse JSON correctly")
	assert.Contains(t, result, "127.0.0.1")
}
