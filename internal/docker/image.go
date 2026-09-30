package docker

import (
	"fmt"
)

// PullImage pulls a Docker image from a registry.
func (c *Client) PullImage(image, tag string) error {
	fullImage := fmt.Sprintf("%s:%s", image, tag)
	cmd := fmt.Sprintf("docker pull %s", fullImage)
	_, err := c.ssh.Execute(cmd)
	return err
}
