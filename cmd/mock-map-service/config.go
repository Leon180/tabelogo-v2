package main

import (
	"log"
	"os"
	"strconv"
	"time"
)

// loadLatencyConfig loads latency configuration from environment variables
func loadLatencyConfig() *LatencyConfig {
	config := &LatencyConfig{
		Enabled:    false,
		MinLatency: 0,
		MaxLatency: 0,
	}

	// Check if latency simulation is enabled
	if enabled := os.Getenv("MOCK_LATENCY_ENABLED"); enabled == "true" {
		config.Enabled = true

		// Load min latency (default: 100ms)
		minMs := 100
		if minStr := os.Getenv("MOCK_LATENCY_MIN_MS"); minStr != "" {
			if parsed, err := strconv.Atoi(minStr); err == nil {
				minMs = parsed
			}
		}
		config.MinLatency = time.Duration(minMs) * time.Millisecond

		// Load max latency (default: 300ms)
		maxMs := 300
		if maxStr := os.Getenv("MOCK_LATENCY_MAX_MS"); maxStr != "" {
			if parsed, err := strconv.Atoi(maxStr); err == nil {
				maxMs = parsed
			}
		}
		config.MaxLatency = time.Duration(maxMs) * time.Millisecond

		log.Printf("⏱️  Latency simulation enabled: %d-%dms", minMs, maxMs)
	} else {
		log.Println("⚡ Latency simulation disabled (fast mode)")
	}

	return config
}
