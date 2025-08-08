package config

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Settings struct {
	LatencyThreshold int `json:"latencyThreshold"`
}

var (
	settings Settings
	mu       sync.RWMutex
)

const settingsFile = "./internal/config/settings.json"

func LoadSettings() error {
	file, err := os.Open(settingsFile)
	if err != nil {
		return fmt.Errorf("failed to open settings file: %w", err)
	}
	defer file.Close()

	var newSettings Settings
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&newSettings); err != nil {
		return fmt.Errorf("failed to decode settings file: %w", err)
	}

	mu.Lock()
	settings = newSettings
	mu.Unlock()
	return nil
}

func GetThreshold() int {
	mu.RLock()
	defer mu.RUnlock()
	return settings.LatencyThreshold
}

func UpdateSettings(newSettings Settings) error {
	file, err := os.Create(settingsFile)
	if err != nil {
		return fmt.Errorf("failed to open settings file for writing: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(newSettings); err != nil {
		return fmt.Errorf("failed to encode settings: %w", err)
	}

	mu.Lock()
	settings = newSettings
	mu.Unlock()

	return nil
}
