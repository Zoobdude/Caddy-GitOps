package main

import (
	"log/slog"
	"net/http"
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

	err := reloadConfiguration()
	if err != nil {
		httpResponseAndLog(w, "Failed to reload configuration: "+err.Error())
		return
	}

	elapsed := time.Since(requestStart)
	slog.Info("✅ Configuration loaded successfully", "elapsed", elapsed)
	w.Write([]byte("Configuration loaded in " + elapsed.String()))
}

func httpResponseAndLog(w http.ResponseWriter, message string) {
	http.Error(w, message, http.StatusInternalServerError)
	slog.Error("❌ " + message)
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
