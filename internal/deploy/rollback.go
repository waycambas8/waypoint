package deploy

import (
	"fmt"

	"github.com/waycambas/waypoint/internal/docker"
	"github.com/waycambas/waypoint/internal/logger"
)

func rollback(dockerClient *docker.Client, container, backupContainer string) error {
	logger.Info("Deployment failed. Rollback started...")
	
	// Remove the failed new container
	logger.Info(fmt.Sprintf("Removing failed container: %s", container))
	dockerClient.RemoveContainer(container)
	
	// Restore backup container
	logger.Info(fmt.Sprintf("Restoring previous container from: %s", backupContainer))
	_, err := dockerClient.RenameContainer(backupContainer, container)
	if err != nil {
		return fmt.Errorf("failed to rename backup container: %w", err)
	}
	
	// Start the restored container
	logger.Info("Starting restored container...")
	err = dockerClient.StartContainer(container)
	if err != nil {
		return fmt.Errorf("failed to start restored container: %w", err)
	}
	
	logger.Success("Previous version restored")
	return nil
}
