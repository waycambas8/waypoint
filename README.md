# waypoint

A framework-agnostic GitHub Action for deploying Docker images to remote servers via SSH.
It handles pulling images, zero-downtime-like deployments with backup containers, health checks, and automatic rollbacks.

## Features
- **Framework Agnostic:** Works with Laravel, Node.js, Go, Python, Java, etc. as long as they are Dockerized.
- **SSH & ProxyJump:** Support for direct SSH connections or via Bastion/Jump hosts (Key and Password auth supported).
- **Auto Server Provisioning:** Automatically installs Docker and Nginx on bare-metal servers if missing.
- **Automated Reverse Proxy:** Automatically generates and tests Nginx `conf.d` virtual hosts to expose your container via a domain.
- **Docker Deployment:** Pulls images, replaces containers safely.
- **Safe Deployment (Rollback):** Backs up the current container and automatically rolls back if the new container fails the health check.
- **Health Checks:** Native HTTP or Docker command-based health checking.
- **Environment Injection:** Pass environment variables securely via GitHub Secrets.

## Usage

```yaml
- name: Deploy
  uses: waycambas8/waypoint@v1.1.1
  with:
    host: ${{ secrets.SERVER_HOST }}
    username: ${{ secrets.SERVER_USERNAME }}
    key: ${{ secrets.SERVER_KEY }}
    # password: ${{ secrets.SERVER_PASSWORD }} # Alternative to key

    # Optional Proxy
    proxy_host: ${{ secrets.PROXY_HOST }}
    proxy_username: ${{ secrets.PROXY_USERNAME }}
    proxy_key: ${{ secrets.PROXY_KEY }}

    # Registry info
    registry: ghcr.io
    registry_username: ${{ github.actor }}
    registry_token: ${{ secrets.GITHUB_TOKEN }}

    # Deployment
    image: ${{ github.repository }}
    tag: latest
    container: my-app

    # Auto Provisioning & Nginx Reverse Proxy (Optional)
    install_deps: "true"
    domain: "api.yourdomain.com"
    proxy_target_port: "3000"

    env_file: ${{ secrets.APP_ENV }}

    ports: |
      3000:3000

    healthcheck: http://localhost:3000/health
```

### SSH Mode (Execute Remote Commands)
```yaml
- name: Execute Remote Command
  uses: waycambas8/waypoint@v1.1.1
  with:
    mode: ssh
    host: ${{ secrets.SERVER_HOST }}
    username: deploy
    key: ${{ secrets.SERVER_KEY }}
    remote_command: |
      docker ps
      uptime
```
