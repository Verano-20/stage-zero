#!/bin/bash

# Enable error handling and logging
set -e
exec > >(tee -a /var/log/user-data.log)
exec 2>&1

echo "=== User-data script started at $(date) ==="

# Set variables from Terraform template
GITHUB_TOKEN="${GITHUB_TOKEN}"
JWT_SECRET="${JWT_SECRET}"
POSTGRES_PASSWORD="${POSTGRES_PASSWORD}"
DB_PASSWORD="${DB_PASSWORD}"

ENVIRONMENT="${ENVIRONMENT}"
DROPLET_NAME="${DROPLET_NAME}"
SERVICE_NAME="${SERVICE_NAME}"
SERVICE_VERSION="${SERVICE_VERSION}"
SERVICE_PORT="${SERVICE_PORT}"
POSTGRES_DB="${POSTGRES_DB}"
POSTGRES_USER="${POSTGRES_USER}"
DB_HOST="${DB_HOST}"
DB_USER="${DB_USER}"
DB_NAME="${DB_NAME}"
DB_PORT="${DB_PORT}"
ENABLE_STDOUT="${ENABLE_STDOUT}"
ENABLE_OTLP="${ENABLE_OTLP}"
OTLP_ENDPOINT="${OTLP_ENDPOINT}"
OTLP_INSECURE="${OTLP_INSECURE}"
METRIC_INTERVAL="${METRIC_INTERVAL}"

echo "Variables set: 
ENVIRONMENT=$ENVIRONMENT,
DROPLET_NAME=$DROPLET_NAME,
SERVICE_NAME=$SERVICE_NAME,
SERVICE_VERSION=$SERVICE_VERSION,
SERVICE_PORT=$SERVICE_PORT,
POSTGRES_DB=$POSTGRES_DB,
POSTGRES_USER=$POSTGRES_USER,
DB_HOST=$DB_HOST,
DB_USER=$DB_USER,
DB_NAME=$DB_NAME,
DB_PORT=$DB_PORT,
ENABLE_STDOUT=$ENABLE_STDOUT,
ENABLE_OTLP=$ENABLE_OTLP,
OTLP_ENDPOINT=$OTLP_ENDPOINT,
OTLP_INSECURE=$OTLP_INSECURE,
METRIC_INTERVAL=$METRIC_INTERVAL"

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
WORK_DIR="/opt/stage-zero"
mkdir -p $WORK_DIR
cd $WORK_DIR
echo "Working directory: $(pwd)"

# Create environment file for deployment
echo "=== Creating environment file ==="
cat > .env.docker.deployment << EOF
# Database Configuration
POSTGRES_DB=${POSTGRES_DB}
POSTGRES_USER=${POSTGRES_USER}
POSTGRES_PASSWORD=${POSTGRES_PASSWORD}

# Application Configuration
SERVICE_NAME=${SERVICE_NAME}
SERVICE_VERSION=${SERVICE_VERSION}
SERVICE_PORT=${SERVICE_PORT}
ENVIRONMENT=${ENVIRONMENT}

# Database Connection
DB_HOST=${DB_HOST}
DB_USER=${DB_USER}
DB_PASSWORD=${DB_PASSWORD}
DB_NAME=${DB_NAME}
DB_PORT=${DB_PORT}

# JWT Configuration
JWT_SECRET=${JWT_SECRET}

# Telemetry Configuration
ENABLE_STDOUT=${ENABLE_STDOUT}
ENABLE_OTLP=${ENABLE_OTLP}
OTLP_ENDPOINT=${OTLP_ENDPOINT}
OTLP_INSECURE=${OTLP_INSECURE}
METRIC_INTERVAL=${METRIC_INTERVAL}
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

# Download the update script
echo "=== Downloading container deployment script ==="
curl -o deploy-containers.sh https://raw.githubusercontent.com/Verano-20/stage-zero/main/scripts/deploy-containers.sh
chmod +x deploy-containers.sh
echo "Update script downloaded and made executable"

# Infrastructure setup complete
echo "=== Infrastructure setup completed ==="
echo "Docker and Docker Compose installed"
echo "Application directory created: $WORK_DIR"
echo "Configuration files downloaded"
echo "Environment file created"
echo "Container deployment script downloaded and made executable"
echo ""
echo "Next step: Run deploy-containers.sh to deploy containers"

echo "Infrastructure setup completed at $(date)" >> /var/log/setup.log
echo "Ready for application deployment" >> /var/log/setup.log

# Create a completion marker file for the GitHub Actions workflow
touch /var/log/user-data-complete
echo "User-data script completed at $(date)" > /var/log/user-data-complete

echo "=== User-data script completed successfully ==="
