package ping

import (
	"testing"
)

func TestCheckHost_Valid(t *testing.T) {
	result := CheckHost("1.1.1.1", 80)

	if result.Error != nil {
		t.Errorf("Expected host to be UP, got error: %v", result.Error)
	}
	if result.Latency <= 0 {
		t.Errorf("Expected latency > 0, got %v", result.Latency)
	}
}

func TestCheckHost_Invalid(t *testing.T) {
	result := CheckHost("256.256.256.256", 80)

	if result.Error == nil {
		t.Errorf("Expected error for invalid host, got none")
	}
}
