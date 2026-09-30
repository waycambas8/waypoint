package ssh

import (
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// parsePrivateKey parses a private key and returns an ssh.AuthMethod.
func parsePrivateKey(privateKey string) (ssh.AuthMethod, error) {
	// If the private key is a file path (starts with /), read it
	keyData := []byte(privateKey)
	if _, err := os.Stat(privateKey); err == nil {
		data, err := os.ReadFile(privateKey)
		if err == nil {
			keyData = data
		}
	}

	signer, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return ssh.PublicKeys(signer), nil
}
