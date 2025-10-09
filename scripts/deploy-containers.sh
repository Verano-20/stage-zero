#!/bin/bash

# Container deployment script for Stage Zero.
# Pulls all docker images. On first run, start all services. On subsequent runs, restart app containers only.

# Services to restart on update
SERVICES_TO_UPDATE="app migrate"

set -e

echo "=== Starting container deployment at $(date) ==="

# Change to container directory
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

# Determine if this is first deployment or update
if docker-compose -f docker-compose.deployment.yml ps | grep -q "Up"; then
    DEPLOYMENT_TYPE="update"
    echo "=== Updating existing deployment ==="
    echo "Stopping services [$SERVICES_TO_UPDATE]..."
    docker-compose -f docker-compose.deployment.yml stop $SERVICES_TO_UPDATE
    
    echo "Starting updated services [$SERVICES_TO_UPDATE]..."
    docker-compose -f docker-compose.deployment.yml up -d --no-deps $SERVICES_TO_UPDATE
else
    DEPLOYMENT_TYPE="initial"
    echo "=== First deployment - starting all services ==="
    echo "Starting all services..."
    docker-compose -f docker-compose.deployment.yml up -d
fi

# Output deployment type for GitHub Actions
echo "DEPLOYMENT_TYPE=$DEPLOYMENT_TYPE"

# Wait for all services to be ready
echo "=== Waiting for all services to start ==="
timeout=120  # 120 seconds
counter=0

while [ $counter -lt $timeout ]; do
    # Check for not ready services
    NOT_READY_SERVICES=$(docker-compose -f docker-compose.deployment.yml ps | grep -E "(Exit|Restarting|Created|Starting)" | wc -l)
    
    if [ "$NOT_READY_SERVICES" -eq 0 ]; then
        echo "✅ All services are ready!"
        break
    fi

    echo "Waiting for services... ($counter/$timeout) - $NOT_READY_SERVICES services not ready"
    sleep 5
    counter=$((counter + 5))
done

if [ $counter -ge $timeout ]; then
    echo "❌ Services failed to start within ${timeout} seconds"
    echo "Service status:"
    docker-compose -f docker-compose.deployment.yml ps
    
    # Create failure marker for GitHub Actions workflow
    touch /var/log/deploy-containers-failed
    echo "Container deployment failed at $(date)" > /var/log/deploy-containers-failed
    echo "Services failed to start within ${timeout} seconds" >> /var/log/deploy-containers-failed
    echo "DEPLOYMENT_TYPE=failed" >> /var/log/deploy-containers-failed
    
    exit 1
fi

# Show final service status
echo "=== Service status after successful deployment ==="
docker-compose -f docker-compose.deployment.yml ps

# Completion
echo "=== Container deployment completed successfully at $(date) ==="
echo "Container deployment completed successfully at $(date)" >> /var/log/deploy-output.log

# Create completion marker for GitHub Actions workflow
touch /var/log/deploy-containers-complete
echo "Container deployment completed successfully at $(date)" > /var/log/deploy-containers-complete
echo "DEPLOYMENT_TYPE=$DEPLOYMENT_TYPE" >> /var/log/deploy-containers-complete
