package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)

	slog.Info("🚀 Starting Caddy GitOps...")
	// Start server in a goroutine
	go startServer()

	// Start long-running operation in another goroutine
	go healthCheck()

	// Keep main process alive and handle graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	slog.Info("🎯 All services running. Press Ctrl+C to exit.")
	<-c // Block until signal is received

	slog.Info("🛑 Shutting down all services...")

}
