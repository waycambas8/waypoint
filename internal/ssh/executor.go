package ssh

import (
	"bytes"
	"fmt"

	"github.com/waycambas/waypoint/internal/logger"
)

// Execute runs a command on the remote host and returns the output.
func (c *Client) Execute(command string) (string, error) {
	return c.executeInternal(command, false)
}

// ExecuteSilent runs a command but does not log the command string (useful for secrets).
func (c *Client) ExecuteSilent(command string) (string, error) {
	return c.executeInternal(command, true)
}

func (c *Client) executeInternal(command string, silent bool) (string, error) {
	session, err := c.client.NewSession()
	if err != nil {
		return "", fmt.Errorf("failed to create session: %w", err)
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf

	if !silent {
		logger.Info(fmt.Sprintf("Executing command: %s", command))
	} else {
		logger.Info("Executing silent command (***)")
	}

	err = session.Run(command)

	stdout := stdoutBuf.String()
	stderr := stderrBuf.String()

	if len(stdout) > 0 {
		logger.Info(fmt.Sprintf("Stdout:\n%s", stdout))
	}
	if len(stderr) > 0 {
		logger.Error(fmt.Sprintf("Stderr:\n%s", stderr))
	}

	if err != nil {
		return stdout, fmt.Errorf("command execution failed: %w (stderr: %s)", err, stderr)
	}

	return stdout, nil
}
