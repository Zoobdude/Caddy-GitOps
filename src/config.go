package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"
)

var reloadMutex sync.Mutex

func reloadConfiguration() error {
	if !reloadMutex.TryLock() {
		slog.Info("Config reload in prog...")
		return fmt.Errorf("configuration reload already in progress")
	}
	defer reloadMutex.Unlock()

	slog.Info("🔄 Reloading Caddy configuration...")
	gitRepo := os.Getenv("GIT_REPO")
	if gitRepo == "" {
		return fmt.Errorf("GIT_REPO environment variable not set")
	}

	authType := os.Getenv("AUTHENTICATION_TYPE")
	if authType == "" {
		return fmt.Errorf("AUTHENTICATION_TYPE environment variable not set")
	}

	caddyAPI := os.Getenv("CADDY_API")
	if caddyAPI == "" {
		return fmt.Errorf("CADDY_API environment variable not set")
	}

	configFilePath := os.Getenv("CONFIG_FILE_PATH")
	if configFilePath == "" {
		return fmt.Errorf("CONFIG_FILE_PATH environment variable not set")
	}

	configFileType := os.Getenv("CONFIG_FILE_TYPE")
	if configFileType == "" {
		return fmt.Errorf("CONFIG_FILE_TYPE environment variable not set")
	}

	cloneConfigRepository(gitRepo, authType)

	configContent := readConfigFile(configFilePath)

	err := loadConfig(caddyAPI, configContent, configFileType)
	if err != nil {
		return err
	}
	return nil
}

func loadConfig(caddyURL string, configValue string, configType string) error {
	slog.Info("📨 Sending configuration to Caddy...")
	slog.Debug("🔍 Caddy configuration details", "url", caddyURL, "config", configValue, "type", configType)

	configData := bytes.NewBuffer([]byte(configValue))

	resp, err := http.Post(caddyURL+"/load", "application/"+configType, configData)
	if err != nil {
		return fmt.Errorf("Failed to load Caddy configuration: %w", err)
	}
	defer resp.Body.Close()

	// Try to read the response body
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return fmt.Errorf("failed to read Caddy API response: %w", err)
	}

	if (len(body) == 0) && resp.StatusCode == http.StatusOK {
		return nil
	}

	slog.Error("❌ Failed to load Caddy configuration", "status", resp.Status, "caddyResponse", string(body))
	return fmt.Errorf("failed to load Caddy configuration: %s", string(body))
}

func readConfigFile(filePath string) string {
	fullPath := "/tmp/Caddy-GitOps/" + filePath
	data, err := os.ReadFile(fullPath)

	if err != nil {
		slog.Error("❌ Failed to read Caddyfile", "error", err)
		return ""
	}
	return string(data)
}

func checkToken(r *http.Request) bool {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		slog.Warn("🚫 Unauthorised access, no token provided")
		return false
	}

	if authHeader != "your-secret-token" {
		slog.Warn("🚫 Unauthorised access", "providedToken", authHeader)
		return false
	}

	slog.Info("🔑 Request authorised")
	return true
}

func cloneConfigRepository(gitURL string, authType string) error {
	slog.Info("📩 Fetching Caddyfile from Git repository...")

	tmpDir := "/tmp/Caddy-GitOps"
	os.RemoveAll(tmpDir) // Clean up any existing temporary directory
	err := os.Mkdir(tmpDir, 0755)
	if err != nil {
		slog.Error("❌ Failed to create temporary clone directory", "error", err)
		return err
	}

	slog.Debug("🔗 Cloning Git repository", "url", gitURL)

	authType = strings.ToLower(authType)
	switch authType {
	case "ssh":
		slog.Debug("🔑 Using SSH authentication for Git")

		// Implement SSH cloning logic here

	case "https":
		slog.Debug("🔑 Using HTTPS authentication for Git")

		// Split the URL to remove protocol for token insertion
		urlParts := strings.Split(gitURL, "://")
		if len(urlParts) != 2 {
			return fmt.Errorf("invalid URL format: %s", gitURL)
		}

		protocol := urlParts[0] // "https"
		repoPath := urlParts[1] // "github.com/user/repo.git"

		username := os.Getenv("HTTPS_AUTHENTICATION_TOKEN")
		token := os.Getenv("HTTPS_AUTHENTICATION_TOKEN")

		gitRepo := fmt.Sprintf("%s://%s:%s@%s", protocol, username, token, repoPath)

		cmd := exec.Command("git", "clone", "--depth", "1", gitRepo, tmpDir)
		data, err := cmd.CombinedOutput()

		if err != nil {
			slog.Error("❌ Failed to clone repository via HTTPS", "error", err, "output", string(data))
			return fmt.Errorf("failed to clone repository: %w", err)
		}

		slog.Info("✅ Successfully cloned repository via HTTPS authentication")

	default:
		slog.Debug("❌🔑 Using no authentication for Git")

		cmd := exec.Command("git", "clone", "--depth", "1", gitURL, tmpDir)
		data, err := cmd.CombinedOutput()

		if err != nil {
			slog.Error("❌ Failed to clone repository via HTTPS", "error", err, "output", string(data))
			return fmt.Errorf("failed to clone repository: %w", err)
		}

		slog.Info("✅ Successfully cloned repository without authentication")
	}
	return nil
}
