package deploy

import (
	"fmt"
	"strings"
	"time"

	"github.com/waycambas/waypoint/internal/logger"
	"github.com/waycambas/waypoint/internal/ssh"
)

// checkHealth performs a health check on the deployed container.
func checkHealth(client *ssh.Client, container, healthCheck string) error {
	if healthCheck == "" {
		return nil // No health check configured
	}

	logger.Info(fmt.Sprintf("Running health check: %s", healthCheck))

	maxRetries := 10
	delay := 3 * time.Second

	for i := 0; i < maxRetries; i++ {
		var cmd string
		if strings.HasPrefix(healthCheck, "http://") || strings.HasPrefix(healthCheck, "https://") {
			cmd = fmt.Sprintf("curl -f -s %s > /dev/null", healthCheck)
		} else {
			// Interpret as command inside container
			cmd = fmt.Sprintf("docker exec %s %s", container, healthCheck)
		}

		_, err := client.Execute(cmd)
		if err == nil {
			logger.Success("Health check passed")
			return nil
		}

		logger.Info(fmt.Sprintf("Health check attempt %d failed, retrying in %v...", i+1, delay))
		time.Sleep(delay)
	}

	return fmt.Errorf("health check failed after %d attempts", maxRetries)
}
