package ssh

import (
	"fmt"
	"net"
	"time"

	"github.com/waycambas/waypoint/internal/config"
	"golang.org/x/crypto/ssh"
)

// Client wraps an ssh.Client to provide a higher-level interface.
type Client struct {
	client *ssh.Client
}

// Connect establishes an SSH connection to a target host, optionally via a proxy.
func Connect(cfg *config.Config) (*Client, error) {
	targetAuth, err := GetAuthMethod(cfg.Key, cfg.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to get target auth method: %w", err)
	}

	targetConfig := &ssh.ClientConfig{
		User: cfg.Username,
		Auth: []ssh.AuthMethod{targetAuth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // In production, we should verify the host key
		Timeout: 30 * time.Second,
	}

	targetAddr := net.JoinHostPort(cfg.Host, cfg.Port)

	// Direct connection
	if cfg.ProxyHost == "" {
		c, err := ssh.Dial("tcp", targetAddr, targetConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to target host: %w", err)
		}
		return &Client{client: c}, nil
	}

	// Connect via proxy
	proxyAuth, err := GetAuthMethod(cfg.ProxyKey, cfg.ProxyPassword)
	if err != nil {
		return nil, fmt.Errorf("failed to get proxy auth method: %w", err)
	}

	proxyConfig := &ssh.ClientConfig{
		User: cfg.ProxyUsername,
		Auth: []ssh.AuthMethod{proxyAuth},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout: 30 * time.Second,
	}

	proxyAddr := net.JoinHostPort(cfg.ProxyHost, cfg.ProxyPort)
	proxyClient, err := ssh.Dial("tcp", proxyAddr, proxyConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to proxy host: %w", err)
	}

	netConn, err := proxyClient.Dial("tcp", targetAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to dial target host from proxy: %w", err)
	}

	conn, chans, reqs, err := ssh.NewClientConn(netConn, targetAddr, targetConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to negotiate SSH connection to target host: %w", err)
	}

	targetClient := ssh.NewClient(conn, chans, reqs)
	return &Client{client: targetClient}, nil
}

// Close closes the underlying SSH client.
func (c *Client) Close() error {
	if c.client != nil {
		return c.client.Close()
	}
	return nil
}
