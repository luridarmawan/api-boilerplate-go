# Deployment Guide

This guide covers deployment strategies for the API across different environments using CI/CD pipelines.

## Table of Contents

- [Overview](#overview)
- [Environment Configuration](#environment-configuration)
- [Build Process](#build-process)
- [CI/CD Pipelines](#cicd-pipelines)
  - [GitHub Actions](#github-actions)
  - [GitLab CI/CD](#gitlab-cicd)
- [Docker Deployment](#docker-deployment)
- [Manual Deployment](#manual-deployment)
- [Environment Variables](#environment-variables)
- [Troubleshooting](#troubleshooting)

## Overview

The API supports multiple deployment environments:
- **Development**: Local development with hot reload
- **Staging**: Pre-production testing environment
- **Production**: Live production environment

Each environment has its own configuration and deployment strategy.

## Environment Configuration

### Environment Files

The project includes environment-specific configuration files:

```
.env.dev      # Development environment
.env.staging  # Staging environment  
.env.prod     # Production environment
.env.example  # Template for new environments
```

### Key Configuration Variables

| Variable | Description | Example |
|----------|-------------|---------|
| `API_NAME` | Application name | "API Hub" |
| `API_DESCRIPTION` | API description | "Global API Hub" |
| `BASEURL` | API base URL | "api.your_domain.id" |
| `SERVER_PORT` | Internal server port | "3000" |
| `DB_HOST` | Database host | "database_host" |
| `DB_NAME` | Database name | "api_hub_production" |

## Build Process

### Local Build Scripts

#### Linux/macOS
```bash
# Build for development
./scripts/build.sh dev

# Build for staging
./scripts/build.sh staging

# Build for production
./scripts/build.sh prod
```

#### Windows
```cmd
REM Build for development
scripts\build.bat dev

REM Build for staging
scripts\build.bat staging

REM Build for production
scripts\build.bat prod
```

### Build Steps

1. **Environment Setup**: Copy environment-specific `.env` file
2. **Swagger Generation**: Generate API documentation with `swag init`
3. **Go Build**: Compile the application binary
4. **Output**: Binary and documentation ready for deployment

## CI/CD Pipelines

### GitHub Actions

Create `.github/workflows/deploy.yml`:

```yaml
name: Deploy API

on:
  push:
    branches: [ main, develop, staging ]
  pull_request:
    branches: [ main ]

env:
  GO_VERSION: '1.21'
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Install dependencies
      run: go mod download
    
    - name: Install swag
      run: go install github.com/swaggo/swag/cmd/swag@latest
    
    - name: Generate Swagger docs
      run: swag init -g cmd/api/main.go -o docs
    
    - name: Run tests
      run: go test -v ./...
    
    - name: Build application
      run: go build -v -o apiserver cmd/api/main.go

  deploy-staging:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/staging'
    environment: staging
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Install swag
      run: go install github.com/swaggo/swag/cmd/swag@latest
    
    - name: Build for staging
      run: |
        chmod +x scripts/build.sh
        ./scripts/build.sh staging
    
    - name: Deploy to staging
      run: |
        # Add your staging deployment commands here
        echo "Deploying to staging environment"
        # Example: rsync, scp, or container deployment
    
  deploy-production:
    needs: test
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    environment: production
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: ${{ env.GO_VERSION }}
    
    - name: Install swag
      run: go install github.com/swaggo/swag/cmd/swag@latest
    
    - name: Build for production
      run: |
        chmod +x scripts/build.sh
        ./scripts/build.sh prod
    
    - name: Build Docker image
      run: |
        docker build -t ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest .
        docker build -t ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }} .
    
    - name: Log in to Container Registry
      uses: docker/login-action@v2
      with:
        registry: ${{ env.REGISTRY }}
        username: ${{ github.actor }}
        password: ${{ secrets.GITHUB_TOKEN }}
    
    - name: Push Docker image
      run: |
        docker push ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:latest
        docker push ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}:${{ github.sha }}
    
    - name: Deploy to production
      run: |
        # Add your production deployment commands here
        echo "Deploying to production environment"
        # Example: kubectl, docker-compose, or cloud deployment
```

### GitLab CI/CD

Create `.gitlab-ci.yml`:

```yaml
stages:
  - test
  - build
  - deploy

variables:
  GO_VERSION: "1.21"
  DOCKER_DRIVER: overlay2
  DOCKER_TLS_CERTDIR: "/certs"

# Cache Go modules
cache:
  paths:
    - .cache/go-build/
    - .cache/go-mod/

before_script:
  - mkdir -p .cache/go-build .cache/go-mod
  - export GOPATH="$CI_PROJECT_DIR/.cache/go-mod"
  - export GOCACHE="$CI_PROJECT_DIR/.cache/go-build"

test:
  stage: test
  image: golang:${GO_VERSION}-alpine
  before_script:
    - apk add --no-cache git
    - go install github.com/swaggo/swag/cmd/swag@latest
  script:
    - go mod download
    - swag init -g cmd/api/main.go -o docs
    - go test -v ./...
    - go build -v -o apiserver cmd/api/main.go
  artifacts:
    paths:
      - apiserver
      - docs/
    expire_in: 1 hour

build:staging:
  stage: build
  image: golang:${GO_VERSION}-alpine
  before_script:
    - apk add --no-cache git
    - go install github.com/swaggo/swag/cmd/swag@latest
  script:
    - chmod +x scripts/build.sh
    - ./scripts/build.sh staging
  artifacts:
    paths:
      - bin/apiserver
      - docs/
    expire_in: 1 day
  only:
    - staging

build:production:
  stage: build
  image: golang:${GO_VERSION}-alpine
  before_script:
    - apk add --no-cache git
    - go install github.com/swaggo/swag/cmd/swag@latest
  script:
    - chmod +x scripts/build.sh
    - ./scripts/build.sh prod
  artifacts:
    paths:
      - bin/apiserver
      - docs/
    expire_in: 1 week
  only:
    - main

# Docker build for production
docker:build:
  stage: build
  image: docker:latest
  services:
    - docker:dind
  before_script:
    - docker login -u $CI_REGISTRY_USER -p $CI_REGISTRY_PASSWORD $CI_REGISTRY
  script:
    - docker build -t $CI_REGISTRY_IMAGE:latest .
    - docker build -t $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA .
    - docker push $CI_REGISTRY_IMAGE:latest
    - docker push $CI_REGISTRY_IMAGE:$CI_COMMIT_SHA
  only:
    - main

deploy:staging:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache openssh-client rsync
    - eval $(ssh-agent -s)
    - echo "$STAGING_SSH_PRIVATE_KEY" | tr -d '\r' | ssh-add -
    - mkdir -p ~/.ssh
    - chmod 700 ~/.ssh
    - ssh-keyscan $STAGING_HOST >> ~/.ssh/known_hosts
  script:
    - rsync -avz --delete bin/ $STAGING_USER@$STAGING_HOST:$STAGING_PATH/
    - rsync -avz --delete docs/ $STAGING_USER@$STAGING_HOST:$STAGING_PATH/docs/
    - ssh $STAGING_USER@$STAGING_HOST "cd $STAGING_PATH && ./restart-staging.sh"
  environment:
    name: staging
    url: https://api=staging.your_domain.id
  only:
    - staging

deploy:production:
  stage: deploy
  image: alpine:latest
  before_script:
    - apk add --no-cache openssh-client rsync
    - eval $(ssh-agent -s)
    - echo "$PRODUCTION_SSH_PRIVATE_KEY" | tr -d '\r' | ssh-add -
    - mkdir -p ~/.ssh
    - chmod 700 ~/.ssh
    - ssh-keyscan $PRODUCTION_HOST >> ~/.ssh/known_hosts
  script:
    - rsync -avz --delete bin/ $PRODUCTION_USER@$PRODUCTION_HOST:$PRODUCTION_PATH/
    - rsync -avz --delete docs/ $PRODUCTION_USER@$PRODUCTION_HOST:$PRODUCTION_PATH/docs/
    - ssh $PRODUCTION_USER@$PRODUCTION_HOST "cd $PRODUCTION_PATH && ./restart-production.sh"
  environment:
    name: production
    url: https://api.your_domain.id
  when: manual
  only:
    - main
```

## Docker Deployment

### Development with Docker Compose

```bash
# Start development environment
docker-compose up -d

# View logs
docker-compose logs -f api

# Stop environment
docker-compose down
```

### Production Docker Deployment

```bash
# Build production image
docker build -t api-hub:latest .

# Run production container
docker run -d \
  --name api-hub \
  -p 3000:3000 \
  --env-file .env.prod \
  api-hub:latest

# Or with docker-compose
docker-compose -f docker-compose.prod.yml up -d
```

### Kubernetes Deployment

Create `k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-hub
  labels:
    app: api-hub
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-hub
  template:
    metadata:
      labels:
        app: api-hub
    spec:
      containers:
      - name: api-hub
        image: ghcr.io/your-org/api-hub:latest
        ports:
        - containerPort: 3000
        env:
        - name: BASEURL
          value: "api.your_domain.id"
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: db-host
        - name: DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: api-secrets
              key: db-password
        livenessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 3000
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: api-hub-service
spec:
  selector:
    app: api-hub
  ports:
    - protocol: TCP
      port: 80
      targetPort: 3000
  type: LoadBalancer
```

## Manual Deployment

### Server Setup

1. **Install Dependencies**:
   ```bash
   # Install Go
   wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
   sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
   export PATH=$PATH:/usr/local/go/bin
   
   # Install swag
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Deploy Application**:
   ```bash
   # Clone repository
   git clone https://github.com/your-org/api-hub.git
   cd api-hub
   
   # Build for production
   ./scripts/build.sh prod
   
   # Create systemd service
   sudo cp scripts/api-hub.service /etc/systemd/system/
   sudo systemctl enable api-hub
   sudo systemctl start api-hub

   # if create systemd service in user
   # location: ~/.config/systemd/user/
   mkdir -p ~/.config/systemd/user
   # copy scripts/ap-server.service to this location

   ```

### Systemd Service

Create `scripts/api-hub.service`:

```ini
[Unit]
Description=API Hub Service
After=network.target

[Service]
Type=simple
User=apiuser
WorkingDirectory=/opt/api-hub
ExecStart=/opt/api-hub/bin/apiserver
Restart=always
RestartSec=5
Environment=PATH=/usr/local/go/bin:/usr/bin:/bin
EnvironmentFile=/opt/api-hub/.env

[Install]
WantedBy=multi-user.target
```

## Environment Variables

### Required Variables

```bash
# Application
API_NAME="API Hub"
API_DESCRIPTION="Global API Hub"
API_VERSION="0.0.1"
BASEURL="api.your_domain.id"
SERVER_PORT=3000

# Database
DB_HOST="database_host"
DB_PORT=5432
DB_USER="postgres"
DB_PASSWORD="secure_password"
DB_NAME="api_hub_production"
DB_SSLMODE="require"
```

### CI/CD Variables

#### GitHub Secrets
- `STAGING_SSH_PRIVATE_KEY`: SSH key for staging deployment
- `PRODUCTION_SSH_PRIVATE_KEY`: SSH key for production deployment
- `GITHUB_TOKEN`: Automatically provided by GitHub

#### GitLab Variables
- `STAGING_HOST`: Staging server hostname
- `STAGING_USER`: SSH username for staging
- `STAGING_PATH`: Deployment path on staging server
- `STAGING_SSH_PRIVATE_KEY`: SSH private key for staging
- `PRODUCTION_HOST`: Production server hostname
- `PRODUCTION_USER`: SSH username for production
- `PRODUCTION_PATH`: Deployment path on production server
- `PRODUCTION_SSH_PRIVATE_KEY`: SSH private key for production

## Troubleshooting

### Common Issues

1. **Swagger Generation Fails**:
   ```bash
   # Install swag
   go install github.com/swaggo/swag/cmd/swag@latest
   
   # Verify installation
   swag --version
   ```

2. **Environment Variables Not Loading**:
   ```bash
   # Check .env file exists
   ls -la .env*
   
   # Verify environment variables
   printenv | grep API_
   ```

3. **Database Connection Issues**:
   ```bash
   # Test database connection
   psql -h $DB_HOST -U $DB_USER -d $DB_NAME -c "SELECT 1;"
   ```

4. **Port Already in Use**:
   ```bash
   # Find process using port
   lsof -i :3000
   
   # Kill process
   kill -9 <PID>
   ```

### Health Checks

- **Application Health**: `GET /health`
- **API Documentation**: `GET /docs`
- **Version Info**: `GET /version`

### Monitoring

```bash
# Check application logs
journalctl -u api-hub -f

# Check system resources
htop
df -h
free -h

# Check network connections
netstat -tlnp | grep :3000
```

## Security Considerations

1. **Environment Variables**: Never commit sensitive data to version control
2. **SSL/TLS**: Always use HTTPS in production
3. **Database**: Use SSL connections and strong passwords
4. **API Keys**: Rotate API keys regularly
5. **Access Control**: Implement proper authentication and authorization
6. **Firewall**: Configure firewall rules to restrict access

## Backup and Recovery

1. **Database Backups**:
   ```bash
   # Create backup
   pg_dump -h $DB_HOST -U $DB_USER $DB_NAME > backup.sql
   
   # Restore backup
   psql -h $DB_HOST -U $DB_USER $DB_NAME < backup.sql
   ```

2. **Application Backups**:
   ```bash
   # Backup application files
   tar -czf api-hub-backup.tar.gz /opt/api-hub
   ```

## Performance Optimization

1. **Go Build Optimization**:
   ```bash
   # Build with optimizations
   CGO_ENABLED=0 go build -ldflags="-w -s" -o apiserver cmd/api/main.go
   ```

2. **Docker Multi-stage Build**: Use the provided Dockerfile for optimized images

3. **Database Optimization**: Configure connection pooling and indexing

4. **Caching**: Implement Redis for caching frequently accessed data


## TROUBLESHOOT

muncul error: `Error "Failed to connect to bus: No medium found"`:
biasanya terjadi ketika sistem D-Bus user tidak berjalan dengan benar. Berikut solusi untuk memperbaikinya di AlmaLinux:

### SOLUSI

1. Pastikan linger diaktifkan:
```
loginctl enable-linger $(whoami)
```

2. Inisialisasi environment D-Bus:
```
export XDG_RUNTIME_DIR=/run/user/$(id -u)
export DBUS_SESSION_BUS_ADDRESS=unix:path=${XDG_RUNTIME_DIR}/bus
```

3. Mulai service D-Bus secara manual:
```
systemctl --user start dbus.socket
systemctl --user start dbus.service
```

4. Verifikasi koneksi:
```
busctl --user list
```


---

For more information, see:
- [API Documentation](../README.md)
- [Development Guide](./DEVELOPMENT.md)
- [Configuration Reference](./CONFIGURATION.md)