package main

import (
	"fmt"
	"os"

	"github.com/waycambas/waypoint/internal/config"
	"github.com/waycambas/waypoint/internal/deploy"
	"github.com/waycambas/waypoint/internal/logger"
	"github.com/waycambas/waypoint/internal/ssh"
)

func main() {
	logger.Init()
	logger.Info("Starting waypoint...")

	cfg, err := config.LoadFromEnv()
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to load configuration: %v", err))
		os.Exit(1)
	}

	if err := cfg.Validate(); err != nil {
		logger.Error(fmt.Sprintf("Configuration validation failed: %v", err))
		os.Exit(1)
	}

	client, err := ssh.Connect(cfg)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect via SSH: %v", err))
		os.Exit(1)
	}
	defer client.Close()
	logger.Success("SSH connection established")

	if cfg.Mode == "ssh" {
		logger.Info("Executing in SSH mode...")
		_, err := client.Execute(cfg.RemoteCommand)
		if err != nil {
			logger.Error(fmt.Sprintf("Remote command failed: %v", err))
			os.Exit(1)
		}
	} else if cfg.Mode == "deploy" {
		logger.Info("Executing in Deploy mode...")
		err := deploy.Deploy(cfg, client)
		if err != nil {
			logger.Error(fmt.Sprintf("Deployment failed: %v", err))
			os.Exit(1)
		}
	} else {
		logger.Error(fmt.Sprintf("Unknown mode: %s", cfg.Mode))
		os.Exit(1)
	}

	logger.Info("Deployment completed successfully")
}
