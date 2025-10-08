#!/bin/bash

# Update script for Stage Zero deployment
# This script pulls the latest Docker images and restarts the application services

set -e

echo "=== Starting application update at $(date) ==="

# Change to application directory
cd /opt/stage-zero

# Verify we're in the right directory
if [ ! -f "docker-compose.deployment.yml" ]; then
    echo "❌ Error: docker-compose.deployment.yml not found in $(pwd)"
    exit 1
fi

# Check if GitHub token is provided
if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ Error: GITHUB_TOKEN environment variable is required"
    exit 1
fi

echo "✅ Environment verified"

# Login to GitHub Container Registry
echo "=== Logging into GitHub Container Registry ==="
echo "$GITHUB_TOKEN" | docker login ghcr.io -u "$GITHUB_TOKEN" --password-stdin
echo "✅ GitHub Container Registry login successful"

# Pull the latest images
echo "=== Pulling latest container images ==="
docker-compose -f docker-compose.deployment.yml pull
echo "✅ Image pull completed"

# Check current service status
echo "=== Current service status ==="
docker-compose -f docker-compose.deployment.yml ps

# Restart application services only (keep DB and monitoring running)
echo "=== Restarting application services ==="
echo "🔄 Stopping application services..."
docker-compose -f docker-compose.deployment.yml stop app migrate

echo "🚀 Starting updated application services..."
docker-compose -f docker-compose.deployment.yml up -d --no-deps app migrate

# Wait for services to be ready with periodic health checks
echo "=== Waiting for services to start ==="
timeout=60  # 60 seconds
counter=0

while [ $counter -lt $timeout ]; do
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Application is ready!"
        break
    fi
    echo "Waiting for application... ($counter/$timeout)"
    sleep 5
    counter=$((counter + 5))
done

if [ $counter -ge $timeout ]; then
    echo "❌ Application failed to start within ${timeout} seconds"
    echo "Service status:"
    docker-compose -f docker-compose.deployment.yml ps
    echo "Service logs:"
    docker-compose -f docker-compose.deployment.yml logs app
    exit 1
fi

# Show final service status
echo "=== Service status after successful update ==="
docker-compose -f docker-compose.deployment.yml ps

# Get droplet IP for logging
DROPLET_IP=$(curl -s http://169.254.169.254/metadata/v1/interfaces/public/0/ipv4/address 2>/dev/null || echo "unknown")

# Log completion
echo "=== Application update completed successfully at $(date) ==="
echo "Application is available at http://$DROPLET_IP:8080"
echo "Grafana dashboard: http://$DROPLET_IP:3000 (admin/admin)"
echo "Prometheus metrics: http://$DROPLET_IP:9090"

# Log to system log
echo "Application update completed successfully at $(date)" >> /var/log/stage-zero-updates.log

echo "🎉 Update completed successfully!"
