package docker

import (
	"fmt"
	"strings"
)

// RunContainer runs a new Docker container.
func (c *Client) RunContainer(image, tag, container, ports, envFile, command string) error {
	fullImage := fmt.Sprintf("%s:%s", image, tag)

	cmdBuilder := []string{"docker", "run", "-d", "--name", container}

	if ports != "" {
		portMappings := strings.Split(ports, "\n")
		for _, p := range portMappings {
			p = strings.TrimSpace(p)
			if p != "" {
				cmdBuilder = append(cmdBuilder, "-p", p)
			}
		}
	}

	if envFile != "" {
		cmdBuilder = append(cmdBuilder, "--env-file", envFile)
	}

	cmdBuilder = append(cmdBuilder, fullImage)

	if command != "" {
		cmdBuilder = append(cmdBuilder, command)
	}

	cmd := strings.Join(cmdBuilder, " ")
	_, err := c.ssh.Execute(cmd)
	return err
}

// StopContainer stops a running container.
func (c *Client) StopContainer(container string) error {
	_, err := c.ssh.Execute(fmt.Sprintf("docker stop %s", container))
	return err // Often ignoring errors here is fine if it doesn't exist, we'll let caller handle it.
}

// RemoveContainer removes a container.
func (c *Client) RemoveContainer(container string) error {
	_, err := c.ssh.Execute(fmt.Sprintf("docker rm -f %s", container))
	return err
}

// ContainerExists checks if a container exists.
func (c *Client) ContainerExists(container string) (bool, error) {
	out, err := c.ssh.Execute(fmt.Sprintf("docker ps -a -q -f name=^/%s$", container))
	if err != nil {
		return false, err
	}
	return strings.TrimSpace(out) != "", nil
}

// RenameContainer renames a container.
func (c *Client) RenameContainer(oldName, newName string) (string, error) {
	return c.ssh.Execute(fmt.Sprintf("docker rename %s %s", oldName, newName))
}

// StartContainer starts a container.
func (c *Client) StartContainer(container string) error {
	_, err := c.ssh.Execute(fmt.Sprintf("docker start %s", container))
	return err
}
