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

	requestURL := healthCheckEndpoint + "/health?key=" + string(sharedHealthCheckKey)
	slog.Info("🔗 Health check URL: " + requestURL)
	resp, err := http.Get(requestURL)

	if err != nil {
		slog.Error("❌ Health check request failed", "error", err)
		return false
	} else if resp.StatusCode == http.StatusNotFound {
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
	healthStatus := true

	healthyInterval := 30 * time.Second // Default value

	unhealthyInterval := 5 * time.Millisecond // Default value}

	slog.Info("🔧 Health check configuration",
		"healthy_interval", healthyInterval,
		"unhealthy_interval", unhealthyInterval)

	for {

		// Unhealthy
		for !healthStatus {
			if healthy() {
				slog.Warn(("❤️‍🩹 Connection restored"))
				healthStatus = true
			} else {
				slog.Info("🤒 Connection still unhealthy, polling Git for changes...")
				err := reloadConfiguration()
				if err != nil {
					slog.Error("❌ Failed to reload configuration", "error", err)
				} else {
					slog.Info("✅ Configuration reloaded successfully")
				}

				//time.Sleep(unhealthyInterval)
			}
		}

		// Healthy
		if !healthy() {
			slog.Warn("😵 Connection lost ")
			healthStatus = false
		} else {
			slog.Info("✅ Health check successful")
			time.Sleep(healthyInterval)
		}
	}
}
