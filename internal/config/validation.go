package config

import (
	"fmt"
)

// Validate checks if the configuration is valid for the selected mode.
func (c *Config) Validate() error {
	if c.Host == "" {
		return fmt.Errorf("host is required")
	}
	if c.Username == "" {
		return fmt.Errorf("username is required")
	}
	if c.Key == "" && c.Password == "" {
		return fmt.Errorf("either key or password is required")
	}

	if c.Mode == "ssh" {
		if c.RemoteCommand == "" {
			return fmt.Errorf("remote_command is required in ssh mode")
		}
	} else if c.Mode == "deploy" {
		if c.Image == "" {
			return fmt.Errorf("image is required in deploy mode")
		}
		if c.Tag == "" {
			return fmt.Errorf("tag is required in deploy mode")
		}
		if c.Container == "" {
			return fmt.Errorf("container is required in deploy mode")
		}
	} else {
		return fmt.Errorf("invalid mode: %s. Must be 'ssh' or 'deploy'", c.Mode)
	}

	return nil
}
