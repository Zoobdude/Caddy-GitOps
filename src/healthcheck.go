package main

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func healthy() bool {
	slog.Info("🌡️ Performing health check...")
	healthCheckEndpoint := os.Getenv("ENDPOINT_ADDRESS")
	if healthCheckEndpoint == "" {
		slog.Error("❌ Health check endpoint not set")
		return false
	}

	// Generate and set request key
	healthCheckKeyMutex.Lock()
	sharedHealthCheckKey = rand.Text()
	healthCheckKeyMutex.Unlock()

	resp, err := http.Get(healthCheckEndpoint + "/health?key=" + string(sharedHealthCheckKey))

	if err != nil {
		slog.Error("❌ Health check request failed", "error", err)
		return false
	} else if resp.StatusCode != http.StatusNotFound {
		slog.Error("❌ Health check endpoint not found, wrong ENDPOINT_ADDRESS?", "statusCode", resp.StatusCode)
		return false

	} else if resp.StatusCode != http.StatusOK {
		slog.Error("❌ Health check failed", "statusCode", resp.StatusCode)
		return false
	}
	defer resp.Body.Close()

	slog.Info("✅ Health check passed")
	return true
}

func healthCheck() {
	for {
		if !healthy() {
			slog.Warn("🔄 Retrying health check in 3 seconds...")
		} else {
			slog.Info("✅ Health check successful")
		}
		time.Sleep(3 * time.Second)
	}
}
