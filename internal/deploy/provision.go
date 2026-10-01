package deploy

import (
	"fmt"

	"github.com/waycambas/waypoint/internal/logger"
	"github.com/waycambas/waypoint/internal/ssh"
)

// provisionDeps checks and installs Docker and Nginx if they are missing.
func provisionDeps(client *ssh.Client) error {
	logger.Info("Checking dependencies (Docker, Nginx)...")

	script := `
	set -e

	# Install Docker if missing
	if ! command -v docker &> /dev/null; then
		echo "Installing Docker..."
		curl -fsSL https://get.docker.com -o get-docker.sh
		sudo sh get-docker.sh
		sudo usermod -aG docker $USER || true
		sudo systemctl enable docker --now || true
	else
		echo "Docker is already installed."
	fi

	# Install Nginx if missing
	if ! command -v nginx &> /dev/null; then
		echo "Installing Nginx..."
		if command -v apt-get &> /dev/null; then
			sudo apt-get update && sudo DEBIAN_FRONTEND=noninteractive apt-get install -y nginx
		elif command -v yum &> /dev/null; then
			sudo yum install -y epel-release && sudo yum install -y nginx
		fi
		sudo systemctl enable nginx --now || true
	else
		echo "Nginx is already installed."
	fi
	`

	_, err := client.Execute(script)
	if err != nil {
		return fmt.Errorf("failed to provision dependencies: %w", err)
	}
	return nil
}

// setupNginx configures Nginx as a reverse proxy for the deployed container.
func setupNginx(client *ssh.Client, domain, targetPort string) error {
	logger.Info(fmt.Sprintf("Setting up Nginx reverse proxy for %s -> 127.0.0.1:%s", domain, targetPort))

	configContent := fmt.Sprintf(`server {
    listen 80;
    server_name %s;

    location / {
        proxy_pass http://127.0.0.1:%s;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
`, domain, targetPort)

	// Write the configuration to conf.d which is widely supported across distributions
	script := fmt.Sprintf(`
	cat << 'EOF' > /tmp/nginx_%s.conf
%s
EOF
	# Move to conf.d
	sudo mv /tmp/nginx_%s.conf /etc/nginx/conf.d/%s.conf
	
	# Clean up any default configs that might conflict on port 80
	if [ -f /etc/nginx/sites-enabled/default ]; then
		sudo rm -f /etc/nginx/sites-enabled/default
	fi

	sudo nginx -t && sudo systemctl reload nginx
	`, domain, configContent, domain, domain)

	_, err := client.Execute(script)
	if err != nil {
		return fmt.Errorf("failed to setup Nginx configuration: %w", err)
	}

	logger.Success(fmt.Sprintf("Nginx configured successfully for %s", domain))
	return nil
}
