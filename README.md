# waypoint

A framework-agnostic GitHub Action for deploying Docker images to remote servers via SSH.
It handles pulling images, zero-downtime-like deployments with backup containers, health checks, and automatic rollbacks.

## Features
- **Framework Agnostic:** Works with Laravel, Node.js, Go, Python, Java, etc. as long as they are Dockerized.
- **SSH & ProxyJump:** Support for direct SSH connections or via Bastion/Jump hosts.
- **Docker Deployment:** Pulls images, replaces containers safely.
- **Safe Deployment (Rollback):** Backs up the current container and automatically rolls back if the new container fails the health check.
- **Health Checks:** Native HTTP or Docker command-based health checking.
- **Environment Injection:** Pass environment variables securely via GitHub Secrets.

## Usage

```yaml
- name: Deploy
  uses: waycambas/waypoint@v1
  with:
    host: ${{ vars.SERVER_HOST }}
    username: ${{ vars.SERVER_USERNAME }}
    key: ${{ secrets.SERVER_KEY }}

    # Optional Proxy
    proxy_host: ${{ vars.PROXY_HOST }}
    proxy_username: ${{ vars.PROXY_USERNAME }}
    proxy_key: ${{ secrets.PROXY_KEY }}

    # Registry info
    registry: ghcr.io
    registry_username: ${{ github.actor }}
    registry_token: ${{ secrets.GITHUB_TOKEN }}

    # Deployment
    image: ghcr.io/your-org/my-app
    tag: ${{ github.sha }}
    container: my-app

    env_file: ${{ secrets.APP_ENV }}

    ports: |
      3000:3000

    healthcheck: http://localhost:3000/health
```

### SSH Mode (Execute Remote Commands)
```yaml
- name: Execute Remote Command
  uses: waycambas/waypoint@v1
  with:
    mode: ssh
    host: ${{ vars.SERVER_HOST }}
    username: deploy
    key: ${{ secrets.SERVER_KEY }}
    remote_command: |
      docker ps
      uptime
```
