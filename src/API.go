package main

import (
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"
)

func startServer() {
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/health", healthCheckHandler)

	http.ListenAndServe(":2020", nil)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("🌐 Update request received")

	// Check for authorization token
	if !checkToken(r) {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	requestStart := time.Now()

	cloneConfigRepository(os.Getenv("GIT_REPO"), os.Getenv("AUTHENTICATION_TYPE"))

	err := loadConfig("http://localhost:2019", readConfigFile("Caddyfile"), "caddyfile")
	if err != nil {
		http.Error(w, "Failed to load configuration", http.StatusInternalServerError)
		slog.Error("❌ Error loading configuration", "error", err)
		return
	}
	elapsed := time.Since(requestStart)
	slog.Info("✅ Configuration loaded successfully", "elapsed", elapsed)
	w.Write([]byte("Configuration loaded in " + elapsed.String()))
}

var sharedHealthCheckKey string
var healthCheckKeyMutex sync.Mutex

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	timeStamp := r.Header.Get("Request-Timestamp")
	if timeStamp == "" {
	}

	requestHealthCheckKey := r.URL.Query().Get("key")

	healthCheckKeyMutex.Lock()
	storedHealthCheckKey := sharedHealthCheckKey
	healthCheckKeyMutex.Unlock()

	if requestHealthCheckKey != storedHealthCheckKey {
		slog.Warn("🚫 Health check failed, invalid key")
		http.Error(w, "Invalid health check key", http.StatusUnauthorized)
		return
	}

	// implement health check logic here
}
