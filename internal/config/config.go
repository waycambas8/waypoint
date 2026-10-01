package config

import (
	"os"
)

// Config holds the configuration for the deploy action.
type Config struct {
	Mode string

	// SSH Target
	Host     string
	Username string
	Password string
	Key      string
	Port     string

	// SSH Proxy
	ProxyHost     string
	ProxyUsername string
	ProxyPassword string
	ProxyKey      string
	ProxyPort     string

	// Remote Command
	RemoteCommand string

	// Docker
	Image            string
	Tag              string
	Container        string
	Registry         string
	RegistryUsername string
	RegistryToken    string
	EnvFile          string
	Ports            string
	ContainerCommand string
	HealthCheck      string

	// Timeouts
	CommandTimeout string

	// Provisioning
	InstallDeps     string
	Domain          string
	ProxyTargetPort string
}

// LoadFromEnv loads configuration from environment variables set by GitHub Actions.
func LoadFromEnv() (*Config, error) {
	return &Config{
		Mode:             getEnvOrDefault("INPUT_MODE", "deploy"),
		Host:             os.Getenv("INPUT_HOST"),
		Username:         os.Getenv("INPUT_USERNAME"),
		Password:         os.Getenv("INPUT_PASSWORD"),
		Key:              os.Getenv("INPUT_KEY"),
		Port:             getEnvOrDefault("INPUT_PORT", "22"),
		ProxyHost:        os.Getenv("INPUT_PROXY_HOST"),
		ProxyUsername:    os.Getenv("INPUT_PROXY_USERNAME"),
		ProxyPassword:    os.Getenv("INPUT_PROXY_PASSWORD"),
		ProxyKey:         os.Getenv("INPUT_PROXY_KEY"),
		ProxyPort:        getEnvOrDefault("INPUT_PROXY_PORT", "22"),
		RemoteCommand:    os.Getenv("INPUT_REMOTE_COMMAND"),
		Image:            os.Getenv("INPUT_IMAGE"),
		Tag:              os.Getenv("INPUT_TAG"),
		Container:        os.Getenv("INPUT_CONTAINER"),
		Registry:         getEnvOrDefault("INPUT_REGISTRY", "ghcr.io"),
		RegistryUsername: os.Getenv("INPUT_REGISTRY_USERNAME"),
		RegistryToken:    os.Getenv("INPUT_REGISTRY_TOKEN"),
		EnvFile:          os.Getenv("INPUT_ENV_FILE"),
		Ports:            os.Getenv("INPUT_PORTS"),
		ContainerCommand: os.Getenv("INPUT_CONTAINER_COMMAND"),
		HealthCheck:      os.Getenv("INPUT_HEALTHCHECK"),
		CommandTimeout:   getEnvOrDefault("INPUT_COMMAND_TIMEOUT", "30m"),
		InstallDeps:      getEnvOrDefault("INPUT_INSTALL_DEPS", "false"),
		Domain:           os.Getenv("INPUT_DOMAIN"),
		ProxyTargetPort:  os.Getenv("INPUT_PROXY_TARGET_PORT"),
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultValue
}
