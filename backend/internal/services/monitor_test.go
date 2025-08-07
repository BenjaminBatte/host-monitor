package services

import (
	"testing"
)

func TestCheckHost_Success(t *testing.T) {
	service := NewMonitorService([]string{"1.1.1.1"})

	host, latency, err := service.CheckHost("1.1.1.1", 80)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if latency <= 0 {
		t.Errorf("Expected latency > 0, got %v", latency)
	}
	if host != "1.1.1.1" {
		t.Errorf("Expected host to be 1.1.1.1, got %v", host)
	}
}

func TestCheckHost_Failure(t *testing.T) {
	service := NewMonitorService([]string{"invalid.host"})

	_, _, err := service.CheckHost("invalid.host", 80)

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
