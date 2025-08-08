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
	"github.com/BenjaminBatte/host-monitor/internal/services"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type mockMetricsStore struct {
	mu      sync.Mutex
	metrics map[string]*models.HostMetrics
}

var _ services.MetricsReader = (*mockMetricsStore)(nil)

func (m *mockMetricsStore) Get(host string) *models.HostMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.metrics[host]
}

func (m *mockMetricsStore) All() map[string]*models.HostMetrics {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make(map[string]*models.HostMetrics, len(m.metrics))
	for k, v := range m.metrics {
		out[k] = v
	}
	return out
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
