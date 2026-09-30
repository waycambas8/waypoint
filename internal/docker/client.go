package docker

import (
	"fmt"

	"github.com/waycambas/waypoint/internal/ssh"
)

// Client handles executing Docker commands via an SSH client.
type Client struct {
	ssh *ssh.Client
}

// NewClient creates a new Docker client over SSH.
func NewClient(sshClient *ssh.Client) *Client {
	return &Client{
		ssh: sshClient,
	}
}

// Login logs into a Docker registry.
func (c *Client) Login(registry, username, token string) error {
	cmd := fmt.Sprintf("echo '%s' | docker login %s -u '%s' --password-stdin", token, registry, username)
	_, err := c.ssh.ExecuteSilent(cmd)
	return err
}
