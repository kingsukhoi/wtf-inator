package background

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"time"

	"github.com/kingsukhoi/wtf-inator/pkg/conf"
)

func StartHealthCheckWorker(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Perform initial health check
	performHealthCheck()

	for {
		select {
		case <-ctx.Done():
			log.Println("Health check worker shutting down...")
			return
		case <-ticker.C:
			performHealthCheck()
		}
	}
}

func performHealthCheck() {
	// Create a client with timeout to avoid hanging requests
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	config := conf.MustGetConfig()

	healthCheckUrl, err := url.Parse(config.Server.Url)
	if err != nil {
		slog.Error("Failed to parse health check url", "error", err)
		return
	}

	healthCheckUrl.JoinPath(config.Server.HealthCheckPath)

	// Replace with your actual health check endpoint
	resp, err := client.Get(healthCheckUrl.String())
	if err != nil {
		slog.Error("Health check failed: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		slog.Info("Health check passed")
	} else {
		slog.Error("Health check failed", "status", resp.StatusCode)
	}
}
