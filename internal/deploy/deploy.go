package deploy

import (
	"encoding/base64"
	"fmt"

	"github.com/waycambas/waypoint/internal/config"
	"github.com/waycambas/waypoint/internal/docker"
	"github.com/waycambas/waypoint/internal/logger"
	"github.com/waycambas/waypoint/internal/ssh"
)

// Deploy executes the deployment workflow.
func Deploy(cfg *config.Config, client *ssh.Client) error {
	dockerClient := docker.NewClient(client)

	if cfg.RegistryUsername != "" && cfg.RegistryToken != "" {
		logger.Info(fmt.Sprintf("Logging into Docker registry: %s", cfg.Registry))
		if err := dockerClient.Login(cfg.Registry, cfg.RegistryUsername, cfg.RegistryToken); err != nil {
			return fmt.Errorf("failed to login to docker registry: %w", err)
		}
	}

	fullImage := fmt.Sprintf("%s/%s", cfg.Registry, cfg.Image)
	if cfg.Registry == "" || cfg.Registry == "docker.io" {
		fullImage = cfg.Image
	}

	logger.Info(fmt.Sprintf("Pulling image: %s:%s", fullImage, cfg.Tag))
	if err := dockerClient.PullImage(fullImage, cfg.Tag); err != nil {
		return fmt.Errorf("failed to pull image: %w", err)
	}

	exists, err := dockerClient.ContainerExists(cfg.Container)
	if err != nil {
		return fmt.Errorf("failed to check if container exists: %w", err)
	}

	backupContainer := fmt.Sprintf("%s-backup", cfg.Container)
	hasBackup := false

	if exists {
		logger.Info(fmt.Sprintf("Previous container detected: %s", cfg.Container))
		
		// Ensure old backup is removed first
		dockerClient.RemoveContainer(backupContainer)

		logger.Info(fmt.Sprintf("Stopping existing container: %s", cfg.Container))
		if err := dockerClient.StopContainer(cfg.Container); err != nil {
			logger.Error(fmt.Sprintf("Warning: Failed to stop container: %v", err))
		}

		logger.Info("Saving current container state as backup...")
		if _, err := dockerClient.RenameContainer(cfg.Container, backupContainer); err != nil {
			return fmt.Errorf("failed to rename container to backup: %w", err)
		}
		hasBackup = true
	}

	envFilePath := ""
	if cfg.EnvFile != "" {
		logger.Info("Setting up environment file on remote host...")
		envFilePath = fmt.Sprintf("/tmp/.env-%s", cfg.Container)
		
		encodedEnv := base64.StdEncoding.EncodeToString([]byte(cfg.EnvFile))
		createEnvCmd := fmt.Sprintf("echo '%s' | base64 -d > %s && chmod 600 %s", encodedEnv, envFilePath, envFilePath)
		if _, err := client.ExecuteSilent(createEnvCmd); err != nil {
			return fmt.Errorf("failed to create env file on remote host: %w", err)
		}
		defer func() {
			logger.Info("Cleaning up remote environment file...")
			client.Execute(fmt.Sprintf("rm -f %s", envFilePath))
		}()
	}

	logger.Info(fmt.Sprintf("Starting new container: %s", cfg.Container))
	err = dockerClient.RunContainer(fullImage, cfg.Tag, cfg.Container, cfg.Ports, envFilePath, cfg.ContainerCommand)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to run new container: %v", err))
		if hasBackup {
			rollback(dockerClient, cfg.Container, backupContainer)
		}
		return fmt.Errorf("deployment failed during container run")
	}

	// Health Check
	err = checkHealth(client, cfg.Container, cfg.HealthCheck)
	if err != nil {
		logger.Error(fmt.Sprintf("Health check failed: %v", err))
		if hasBackup {
			rollback(dockerClient, cfg.Container, backupContainer)
		}
		return fmt.Errorf("deployment failed during health check")
	}

	// Cleanup
	if hasBackup {
		logger.Info("Cleaning up backup container...")
		dockerClient.RemoveContainer(backupContainer)
	}

	logger.Success("Deployment successful")
	return nil
}
