#!/bin/bash

# Enable error handling and logging
set -e
exec > >(tee -a /var/log/user-data.log)
exec 2>&1

echo "=== User-data script started at $(date) ==="

# Set variables from Terraform template
GITHUB_TOKEN="${github_token}"
GITHUB_USERNAME="${github_username}"
DROPLET_NAME="${droplet_name}"

echo "Variables set: GITHUB_USERNAME=$GITHUB_USERNAME, DROPLET_NAME=$DROPLET_NAME"

# Update system
echo "=== Updating system packages ==="
export DEBIAN_FRONTEND=noninteractive
apt-get update
apt-get upgrade -y -o Dpkg::Options::="--force-confdef" -o Dpkg::Options::="--force-confold"

# Install Docker
echo "=== Installing Docker ==="
curl -fsSL https://get.docker.com -o get-docker.sh
sh get-docker.sh
systemctl enable docker
systemctl start docker
echo "Docker installation completed"

# Install Docker Compose
echo "=== Installing Docker Compose ==="
curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose
echo "Docker Compose installation completed"

# Verify installations
docker --version
docker-compose --version

# Create application directory
echo "=== Creating application directory ==="
mkdir -p /opt/stage-zero
cd /opt/stage-zero
echo "Working directory: $(pwd)"

# Create environment file for deployment
echo "=== Creating environment file ==="
cat > .env.docker.deployment << EOF
# Database Configuration
POSTGRES_DB=stage_zero
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres123

# Application Configuration
SERVICE_NAME=stage-zero-api
SERVICE_VERSION=1.0.0
SERVICE_PORT=8080
ENVIRONMENT=deployment

# Database Connection
DB_HOST=db
DB_USER=postgres
DB_PASSWORD=postgres123
DB_NAME=stage_zero
DB_PORT=5432

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production-32chars-minimum

# Telemetry Configuration
ENABLE_STDOUT=true
ENABLE_OTLP=true
OTLP_ENDPOINT=otel-collector:4317
OTLP_INSECURE=true
METRIC_INTERVAL=30s
EOF

# Download docker-compose.deployment.yml from the repository
echo "=== Downloading configuration files ==="
curl -o docker-compose.deployment.yml https://raw.githubusercontent.com/Verano-20/stage-zero/main/docker-compose.deployment.yml
echo "Downloaded docker-compose.deployment.yml"

# Download configuration files needed by the services
curl -o otel-collector-config.yaml https://raw.githubusercontent.com/Verano-20/stage-zero/main/otel-collector-config.yaml
echo "Downloaded otel-collector-config.yaml"
curl -o prometheus.yml https://raw.githubusercontent.com/Verano-20/stage-zero/main/prometheus.yml
echo "Downloaded prometheus.yml"
curl -o tempo.yaml https://raw.githubusercontent.com/Verano-20/stage-zero/main/tempo.yaml
echo "Downloaded tempo.yaml"
curl -o loki.yaml https://raw.githubusercontent.com/Verano-20/stage-zero/main/loki.yaml
echo "Downloaded loki.yaml"
curl -o promtail.yaml https://raw.githubusercontent.com/Verano-20/stage-zero/main/promtail.yaml
echo "Downloaded promtail.yaml"

# Create grafana directory and download configuration
echo "=== Setting up Grafana configuration ==="
mkdir -p grafana/provisioning/datasources
mkdir -p grafana/provisioning/dashboards
mkdir -p grafana/dashboards

curl -o grafana/provisioning/datasources/datasources.yml https://raw.githubusercontent.com/Verano-20/stage-zero/main/grafana/provisioning/datasources/datasources.yml
curl -o grafana/provisioning/dashboards/dashboards.yml https://raw.githubusercontent.com/Verano-20/stage-zero/main/grafana/provisioning/dashboards/dashboards.yml
curl -o grafana/dashboards/health.json https://raw.githubusercontent.com/Verano-20/stage-zero/main/grafana/dashboards/health.json
curl -o grafana/dashboards/logs.json https://raw.githubusercontent.com/Verano-20/stage-zero/main/grafana/dashboards/logs.json
curl -o grafana/dashboards/metrics.json https://raw.githubusercontent.com/Verano-20/stage-zero/main/grafana/dashboards/metrics.json
curl -o grafana/dashboards/tracing.json https://raw.githubusercontent.com/Verano-20/stage-zero/main/grafana/dashboards/tracing.json
echo "Grafana configuration completed"

# Login to GitHub Container Registry
echo "=== Logging into GitHub Container Registry ==="
if [ -n "$GITHUB_TOKEN" ] && [ -n "$GITHUB_USERNAME" ]; then
    echo "$GITHUB_TOKEN" | docker login ghcr.io -u "$GITHUB_USERNAME" --password-stdin
    echo "GitHub Container Registry login successful"
else
    echo "ERROR: GitHub credentials not provided"
    exit 1
fi

# Pull the latest images
echo "=== Pulling container images ==="
docker-compose -f docker-compose.deployment.yml pull
echo "Image pull completed"

# Start the application
echo "=== Starting application services ==="
docker-compose -f docker-compose.deployment.yml up -d
echo "Services started"

# Wait for services to be healthy
echo "=== Waiting for services to start ==="
sleep 30

# Check if services are running
echo "=== Checking service status ==="
docker-compose -f docker-compose.deployment.yml ps

# Get droplet IP
DROPLET_IP=$(curl -s http://169.254.169.254/metadata/v1/interfaces/public/0/ipv4/address)

# Log completion
echo "=== Deployment completed at $(date) ==="
echo "Docker setup completed at $(date)" >> /var/log/setup.log
echo "Application should be available at http://$DROPLET_IP:8080" >> /var/log/setup.log
echo "Grafana available at http://$DROPLET_IP:3000" >> /var/log/setup.log
echo "Prometheus available at http://$DROPLET_IP:9090" >> /var/log/setup.log

echo "=== User-data script completed successfully ==="
