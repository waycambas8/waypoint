# Team Deploy GitHub Action

## 1. Overview

Membangun GitHub Action internal berbasis **Golang** untuk kebutuhan:

* SSH ke remote server
* SSH melalui Proxy/Jump Host
* Eksekusi command remote
* Deployment aplikasi berbasis Docker
* Login dan pull image dari GHCR
* Deployment multi-container / multi-service
* Environment file injection
* Health check
* Automatic rollback
* Cleanup image/container lama
* Logging deployment yang jelas
* Reusable oleh seluruh project dalam satu team

Action harus bersifat **framework-agnostic**.

Artinya action tidak perlu mengetahui apakah aplikasi menggunakan:

* Laravel
* Node.js
* NestJS
* Go
* Python
* Java
* .NET
* Rust
* atau framework lainnya

Kontrak deployment adalah **Docker Image**.

---

# 2. Target Usage

Target penggunaan:

```yaml
- name: Deploy
  uses: your-org/deploy-action@v1
  with:
    host: ${{ vars.SERVER_HOST }}
    username: ${{ vars.SERVER_USERNAME }}
    key: ${{ secrets.SERVER_KEY }}

    proxy_host: ${{ vars.PROXY_HOST }}
    proxy_username: ${{ vars.PROXY_USERNAME }}
    proxy_key: ${{ secrets.PROXY_KEY }}

    registry: ghcr.io
    registry_username: ${{ vars.GHCR_USERNAME }}
    registry_token: ${{ secrets.GHCR_TOKEN }}

    image: ghcr.io/your-org/my-app
    tag: ${{ github.sha }}

    container: my-app

    env_file: ${{ secrets.APP_ENV }}

    ports: |
      3000:3000

    healthcheck: http://localhost:3000/health
```

Tujuan akhirnya adalah menghilangkan deployment script Bash yang berulang di setiap repository.

---

# 3. Design Principles

## 3.1 Framework Agnostic

Jangan membuat logic:

```text
if Laravel
if Node.js
if Go
if Python
```

Action hanya menangani deployment artifact.

Artifact:

```text
Docker Image
```

Contoh:

```text
ghcr.io/company/laravel-api:v1
ghcr.io/company/node-api:v1
ghcr.io/company/payment-service:v1
```

Semua diperlakukan sama.

---

## 3.2 Separate Build and Deploy

Build dilakukan oleh project:

```text
Source Code
    ↓
GitHub Actions
    ↓
Docker Build
    ↓
GHCR
```

Deployment dilakukan oleh:

```text
deploy-action
    ↓
SSH
    ↓
Docker Pull
    ↓
Docker Run
```

Action tidak bertanggung jawab terhadap proses build aplikasi.

---

## 3.3 Deployment Should Be Safe

Jangan langsung:

```text
stop container
↓
remove container
↓
pull image
```

Karena apabila pull gagal, aplikasi lama sudah terhapus.

Urutan deployment:

```text
1. Connect SSH
2. Authenticate registry
3. Pull new image
4. Validate image
5. Prepare environment
6. Stop old container
7. Start new container
8. Health check
9. Success → cleanup old resources
10. Failure → rollback
```

---

# 4. Architecture

```text
GitHub Actions
       │
       ▼
┌─────────────────────────────┐
│      deploy-action          │
│                             │
│       Go Binary             │
│                             │
│ ┌─────────────────────────┐ │
│ │ Config                  │ │
│ ├─────────────────────────┤ │
│ │ SSH Client              │ │
│ │ Proxy / Jump Host       │ │
│ │ Remote Executor         │ │
│ │ Registry Client         │ │
│ │ Docker Deployer         │ │
│ │ Health Checker          │ │
│ │ Rollback Manager        │ │
│ │ Logger                  │ │
│ └─────────────────────────┘ │
└──────────────┬──────────────┘
               │
               │ SSH / Proxy
               ▼
        ┌───────────────┐
        │ Remote Server │
        │               │
        │ Docker        │
        │ ├── App       │
        │ ├── Worker    │
        │ └── Service   │
        └───────────────┘
```

---

# 5. Technology Stack

## Core

* Golang
* Go Modules
* Linux
* Docker
* GitHub Actions

## SSH

Recommended:

```text
golang.org/x/crypto/ssh
```

Features:

* Private key authentication
* Password authentication
* SSH agent
* Custom SSH port
* Host key verification
* Connection timeout
* Command timeout

## Proxy

Support:

```text
GitHub Runner
      ↓
Proxy / Bastion
      ↓
Target Server
```

Equivalent concept:

```text
ProxyJump
```

---

# 6. Repository Structure

Recommended structure:

```text
deploy-action/
│
├── action.yml
├── Dockerfile
├── README.md
├── LICENSE
├── Makefile
├── go.mod
├── go.sum
│
├── cmd/
│   └── deploy/
│       └── main.go
│
├── internal/
│   │
│   ├── config/
│   │   ├── config.go
│   │   └── validation.go
│   │
│   ├── ssh/
│   │   ├── client.go
│   │   ├── auth.go
│   │   ├── proxy.go
│   │   └── executor.go
│   │
│   ├── registry/
│   │   └── ghcr.go
│   │
│   ├── docker/
│   │   ├── client.go
│   │   ├── image.go
│   │   ├── container.go
│   │   └── cleanup.go
│   │
│   ├── deploy/
│   │   ├── deploy.go
│   │   ├── service.go
│   │   ├── healthcheck.go
│   │   └── rollback.go
│   │
│   └── logger/
│       └── logger.go
│
└── .github/
    └── workflows/
        ├── test.yml
        └── release.yml
```

---

# 7. GitHub Action Interface

`action.yml` menjadi public interface.

Initial inputs:

```yaml
inputs:

  mode:
    description: "Execution mode: ssh or deploy"
    required: false
    default: "deploy"

  host:
    description: "Target SSH host"
    required: true

  username:
    description: "Target SSH username"
    required: true

  key:
    description: "Target SSH private key"
    required: true

  port:
    description: "Target SSH port"
    required: false
    default: "22"

  proxy_host:
    description: "SSH proxy/jump host"
    required: false

  proxy_username:
    description: "SSH proxy username"
    required: false

  proxy_key:
    description: "SSH proxy private key"
    required: false

  proxy_port:
    description: "SSH proxy port"
    required: false
    default: "22"

  command:
    description: "Remote command for ssh mode"
    required: false

  image:
    description: "Docker image"
    required: false

  tag:
    description: "Docker image tag"
    required: false

  container:
    description: "Docker container name"
    required: false

  registry:
    description: "Docker registry"
    required: false
    default: "ghcr.io"

  registry_username:
    description: "Docker registry username"
    required: false

  registry_token:
    description: "Docker registry token"
    required: false

  env_file:
    description: "Application environment file"
    required: false

  ports:
    description: "Docker port mappings"
    required: false

  command:
    description: "Container startup command"
    required: false

  healthcheck:
    description: "Health check URL"
    required: false

  command_timeout:
    description: "Remote command timeout"
    required: false
    default: "30m"
```

Note:

`command` perlu dipisahkan menjadi:

```text
remote_command
container_command
```

agar tidak terjadi conflict.

---

# 8. SSH Mode

V1 harus tetap menyediakan kemampuan SSH-only.

Example:

```yaml
- name: Execute Remote Command
  uses: your-org/deploy-action@v1
  with:
    mode: ssh

    host: ${{ vars.SERVER_HOST }}
    username: deploy
    key: ${{ secrets.SERVER_KEY }}

    command: |
      docker ps
      df -h
      systemctl status nginx
```

Tujuan:

* Menjadi alternatif internal untuk `appleboy/ssh-action`
* Bisa digunakan untuk maintenance
* Bisa digunakan tanpa Docker deployment

---

# 9. Deploy Mode

Example:

```yaml
- name: Deploy
  uses: your-org/deploy-action@v1
  with:
    mode: deploy

    host: ${{ vars.SERVER_HOST }}
    username: deploy
    key: ${{ secrets.SERVER_KEY }}

    image: ghcr.io/company/payment-api
    tag: ${{ github.sha }}

    container: payment-api

    registry_username: ${{ vars.GHCR_USERNAME }}
    registry_token: ${{ secrets.GHCR_TOKEN }}

    ports: |
      3000:3000

    healthcheck: http://localhost:3000/health
```

---

# 10. Deployment Flow

```text
START
  │
  ▼
Validate Configuration
  │
  ▼
Connect SSH
  │
  ▼
Connect Proxy if configured
  │
  ▼
Authenticate GHCR
  │
  ▼
Pull New Image
  │
  ├── FAIL → Abort
  │
  ▼
Inspect Image
  │
  ▼
Save Current Container State
  │
  ▼
Create New Container
  │
  ▼
Start Container
  │
  ▼
Health Check
  │
  ├── FAIL
  │    │
  │    ▼
  │  Stop New Container
  │    │
  │    ▼
  │  Restore Previous Container
  │
  └── SUCCESS
       │
       ▼
   Cleanup Old Image
       │
       ▼
      DONE
```

---

# 11. Environment File

Support:

```yaml
env_file: ${{ secrets.APP_ENV }}
```

Action akan mengirim environment file ke remote server melalui SSH.

Contoh:

```text
GitHub Secret
      │
      ▼
deploy-action
      │
      ▼
SSH
      │
      ▼
/tmp/deploy/.env
```

Security requirements:

* Jangan print isi `.env`
* Jangan print secret
* Permission file harus restricted
* Cleanup temporary file setelah deployment
* Jangan menyimpan environment file di Git

---

# 12. Docker Deployment

Basic command abstraction:

```text
docker login
docker pull
docker inspect
docker stop
docker rm
docker run
docker healthcheck
docker image prune
```

Action tidak boleh mengasumsikan application framework.

---

# 13. Multi-Service Deployment

Support:

```yaml
services:

  - name: nobu
    image: ghcr.io/company/nobu
    tag: v1.0.0
    container: qris-nobu
    ports:
      - "3000:3000"

  - name: transaction
    image: ghcr.io/company/transaction
    tag: v1.0.0
    container: qris-transaction
    ports:
      - "3001:3001"
```

Deployment:

```text
Application
│
├── nobu
│   ├── pull
│   ├── deploy
│   └── health check
│
└── transaction
    ├── pull
    ├── deploy
    └── health check
```

---

# 14. Health Check

V1 support:

```text
HTTP
TCP
Command
```

Example HTTP:

```yaml
healthcheck:
  type: http
  url: http://localhost:3000/health
  timeout: 30s
  retries: 5
```

Example command:

```yaml
healthcheck:
  type: command
  command: docker inspect --format='{{.State.Running}}' payment-api
```

---

# 15. Rollback

Before replacing container:

```text
Current:
payment-api:v1.2.0
```

New:

```text
payment-api:v1.3.0
```

Save:

```text
previous_image
previous_env
previous_ports
previous_command
```

If health check fails:

```text
v1.3.0
   ↓
FAILED
   ↓
remove v1.3.0
   ↓
restore v1.2.0
```

Output:

```text
Deployment failed.
Rollback started...
Previous version restored.
```

---

# 16. Logging

Output should be readable in GitHub Actions.

Example:

```text
========================================
 Team Deploy Action
========================================

[INFO] Validating configuration...
[INFO] Connecting to proxy...
[INFO] Connecting to target server...
[INFO] SSH connection established

[INFO] Docker registry: ghcr.io
[INFO] Pulling image:
       ghcr.io/company/payment-api:v1.3.0

[INFO] Previous container detected:
       payment-api

[INFO] Starting new container...
[INFO] Running health check...

[OK] Health check passed
[OK] Deployment successful

========================================
 Deployment completed
 Version: v1.3.0
 Container: payment-api
========================================
```

Never output:

```text
SSH private key
GHCR token
.env content
password
secret environment variable
```

---

# 17. Security Requirements

Mandatory:

* Never log private key
* Never log registry token
* Never log password
* Never log `.env`
* Never echo GitHub secrets
* Validate SSH host key when configured
* Prefer GitHub OIDC/token mechanisms where applicable
* Temporary files must be removed
* Avoid `set -x` equivalent behavior
* Validate shell arguments before remote execution
* Prevent accidental command injection

---

# 18. Timeout

Support separate timeouts:

```text
connection_timeout
command_timeout
pull_timeout
healthcheck_timeout
deployment_timeout
```

Example:

```yaml
command_timeout: 30m
healthcheck_timeout: 60s
```

This is important for:

```text
docker pull
```

yang dapat membutuhkan waktu lama.

---

# 19. Retry

Optional retry:

```yaml
retry:
  connection: 3
  pull: 2
  healthcheck: 5
```

Example:

```text
SSH connection failed
       ↓
retry 1
       ↓
retry 2
       ↓
retry 3
       ↓
deployment failed
```

---

# 20. Versioning

Use Semantic Versioning:

```text
v1.0.0
v1.1.0
v1.1.1
v2.0.0
```

Recommended GitHub usage:

```yaml
uses: your-org/deploy-action@v1
```

Internally:

```text
v1
 └── latest compatible v1.x.x
```

Do not recommend:

```yaml
uses: your-org/deploy-action@master
```

for production deployment.

---

# 21. Testing

## Unit Test

Test:

* Config parser
* Config validation
* SSH configuration
* Proxy configuration
* Docker command generation
* Health check
* Rollback logic

## Integration Test

Create temporary Docker environment:

```text
GitHub Runner
     ↓
Test SSH Server
     ↓
Docker
     ↓
Test Container
```

Test:

```text
SSH success
SSH failure
Proxy success
Proxy failure
Docker pull success
Docker pull failure
Container start failure
Health check success
Health check failure
Rollback
```

---

# 22. Development Milestones

## Phase 1 — Project Bootstrap

* [ ] Create GitHub repository
* [ ] Initialize Go module
* [ ] Create project structure
* [ ] Create `action.yml`
* [ ] Create Dockerfile
* [ ] Create basic Go application
* [ ] GitHub Action hello-world test

---

## Phase 2 — SSH Engine

* [ ] SSH private key authentication
* [ ] SSH username
* [ ] SSH host
* [ ] SSH port
* [ ] Connection timeout
* [ ] Command execution
* [ ] stdout streaming
* [ ] stderr streaming
* [ ] exit code handling

---

## Phase 3 — Proxy / Jump Host

* [ ] Proxy host
* [ ] Proxy username
* [ ] Proxy key
* [ ] Proxy port
* [ ] Proxy connection
* [ ] Proxy → target SSH connection
* [ ] Connection timeout

---

## Phase 4 — SSH Mode

Implement:

```yaml
mode: ssh
```

Support:

```yaml
command: |
  docker ps
  uptime
```

At this point action can already replace basic internal use cases of `appleboy/ssh-action`.

---

## Phase 5 — Docker Deployment

Implement:

* [ ] Docker login
* [ ] Docker pull
* [ ] Container detection
* [ ] Container stop
* [ ] Container removal
* [ ] Docker run
* [ ] Port mapping
* [ ] Environment file
* [ ] Container command

---

## Phase 6 — GHCR

Implement:

* [ ] Registry configuration
* [ ] Username
* [ ] Token
* [ ] `docker login`
* [ ] Private image pull
* [ ] Authentication error handling

---

## Phase 7 — Health Check

Implement:

* [ ] HTTP health check
* [ ] TCP health check
* [ ] Command health check
* [ ] Retry
* [ ] Timeout
* [ ] Deployment failure handling

---

## Phase 8 — Rollback

Implement:

* [ ] Capture previous container
* [ ] Preserve previous image
* [ ] Start new version
* [ ] Health check
* [ ] Rollback on failure
* [ ] Cleanup after success

---

## Phase 9 — Multi-Service

Implement:

```yaml
services:
  - ...
  - ...
```

Support:

* [ ] Multiple images
* [ ] Multiple containers
* [ ] Multiple ports
* [ ] Individual health checks
* [ ] Deployment order
* [ ] Failure handling

---

## Phase 10 — Production Hardening

* [ ] Security review
* [ ] Secret masking
* [ ] SSH host verification
* [ ] Input validation
* [ ] Timeout handling
* [ ] Retry
* [ ] Structured logs
* [ ] Integration tests
* [ ] Documentation
* [ ] Release pipeline

---

# 23. V1 Scope

V1 should remain focused.

### Must Have

```text
SSH
ProxyJump
GHCR
Docker
Environment
Single container
Health check
Rollback
Logging
```

### Nice to Have

```text
Multi-service
SFTP
Docker Compose
Deployment hooks
Notifications
```

### Do Not Implement Yet

```text
Kubernetes
Blue/Green deployment
Canary deployment
Service mesh
Infrastructure provisioning
Cloud provider management
```

These can be considered after the basic deployment engine is stable.

---

# 24. Example: Laravel

Docker image:

```text
ghcr.io/company/laravel-api:v1.0.0
```

Deployment action:

```yaml
- uses: your-org/deploy-action@v1
  with:
    mode: deploy

    image: ghcr.io/company/laravel-api
    tag: v1.0.0

    container: laravel-api

    ports: |
      8000:8000

    env_file: ${{ secrets.APP_ENV }}

    healthcheck: http://localhost:8000/up
```

Action does not contain Laravel-specific logic.

---

# 25. Example: Node.js

```yaml
- uses: your-org/deploy-action@v1
  with:
    mode: deploy

    image: ghcr.io/company/node-api
    tag: v2.0.0

    container: node-api

    ports: |
      3000:3000

    healthcheck: http://localhost:3000/health
```

---

# 26. Example: Go

```yaml
- uses: your-org/deploy-action@v1
  with:
    mode: deploy

    image: ghcr.io/company/report-service
    tag: v1.5.0

    container: report-service

    ports: |
      8080:8080

    healthcheck: http://localhost:8080/health
```

All three applications use the same deployment engine.

---

# 27. Future Architecture

Possible future structure:

```text
                    Deploy Action
                         │
             ┌───────────┴───────────┐
             │                       │
          SSH Mode               Deploy Mode
             │                       │
             ▼                       ▼
       Remote Command          Deployment Engine
                                     │
                      ┌──────────────┼──────────────┐
                      ▼              ▼              ▼
                    Docker         Registry       Health
                      │              │              │
                      └──────────────┼──────────────┘
                                     ▼
                                  Rollback
```

Potential future support:

```text
Docker
Docker Compose
Podman
SFTP
Remote scripts
Blue/Green
Canary
Multiple servers
Load balancer integration
Deployment notifications
```

---

# 28. Success Criteria

Project dianggap berhasil apabila team dapat mengganti deployment script yang panjang seperti:

```text
SSH
docker login
docker stop
docker rm
docker rmi
docker pull
docker run
health check
```

menjadi:

```yaml
- uses: your-org/deploy-action@v1
  with:
    ...
```

Dan deployment Laravel, Node.js, Go, atau aplikasi lain dapat menggunakan action yang sama tanpa perubahan pada deployment engine.

---

# 29. Final Goal

Target akhir:

```text
Developer
    │
    │ git push / tag
    ▼
GitHub Actions
    │
    ▼
Build Docker Image
    │
    ▼
GHCR
    │
    ▼
your-org/deploy-action
    │
    ├── SSH
    ├── ProxyJump
    ├── GHCR Login
    ├── Docker Pull
    ├── Environment
    ├── Container Deployment
    ├── Health Check
    └── Rollback
    │
    ▼
Production / Development Server
```

Dengan prinsip:

> **Build once, deploy anywhere.**

Deployment action tidak mengetahui framework aplikasi. Ia hanya mengetahui bagaimana mengambil artifact berupa Docker image dan menjalankannya secara aman di remote server.
