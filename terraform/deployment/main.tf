locals {
    # Droplet Configuration
    DROPLET_NAME = "stage-zero-tf"
    REGION = "lon1"
    # Database Configuration
    POSTGRES_DB="stage_zero"
    POSTGRES_USER="postgres"
    # Application Configuration
    SERVICE_NAME="stage-zero-api"
    SERVICE_VERSION="1.0.0"
    SERVICE_PORT="8080"
    ENVIRONMENT="deployment"
    # Database Connection
    DB_HOST="db"
    DB_USER="postgres"
    DB_NAME="stage_zero"
    DB_PORT="5432"
    # Telemetry Configuration
    ENABLE_STDOUT="false"
    ENABLE_OTLP="true"
    OTLP_ENDPOINT="otel-collector:4317"
    OTLP_INSECURE="true"
    METRIC_INTERVAL="30s"
}

resource "digitalocean_droplet" "stage-zero" {
  name   = local.DROPLET_NAME
  image  = "ubuntu-25-04-x64"
  size   = "s-1vcpu-1gb"
  region = local.REGION
  ssh_keys = [data.digitalocean_ssh_key.terraform.id]
  tags   = ["terraform", local.DROPLET_NAME]

  user_data = templatefile("../../scripts/user-data.sh", {
    # Sensitive
    GITHUB_TOKEN = var.github_token
    JWT_SECRET = var.jwt_secret
    POSTGRES_PASSWORD = var.postgres_password
    DB_PASSWORD = var.db_password
    # Not sensitive
    ENVIRONMENT = local.ENVIRONMENT
    DROPLET_NAME = local.DROPLET_NAME
    SERVICE_NAME = local.SERVICE_NAME
    SERVICE_VERSION = local.SERVICE_VERSION
    SERVICE_PORT = local.SERVICE_PORT
    POSTGRES_DB = local.POSTGRES_DB
    POSTGRES_USER = local.POSTGRES_USER
    DB_HOST = local.DB_HOST
    DB_USER = local.DB_USER
    DB_NAME = local.DB_NAME
    DB_PORT = local.DB_PORT
    ENABLE_STDOUT = local.ENABLE_STDOUT
    ENABLE_OTLP = local.ENABLE_OTLP
    OTLP_ENDPOINT = local.OTLP_ENDPOINT
    OTLP_INSECURE = local.OTLP_INSECURE
    METRIC_INTERVAL = local.METRIC_INTERVAL
  })

  lifecycle {
    prevent_destroy = true
    ignore_changes = [user_data]  # Prevent recreation on user_data changes
  }
}

output "droplet_ip" {
  value = digitalocean_droplet.stage-zero.ipv4_address
  description = "The public IP address of the droplet"
}

output "droplet_id" {
  value = digitalocean_droplet.stage-zero.id
  description = "The ID of the droplet"
}

output "droplet_name" {
  value = digitalocean_droplet.stage-zero.name
  description = "The name of the droplet"
}

output "application_url" {
  value = "http://${digitalocean_droplet.stage-zero.ipv4_address}:8080"
  description = "URL to access the Stage Zero application"
}

output "grafana_url" {
  value = "http://${digitalocean_droplet.stage-zero.ipv4_address}:3000"
  description = "URL to access Grafana dashboard (admin/admin)"
}

output "prometheus_url" {
  value = "http://${digitalocean_droplet.stage-zero.ipv4_address}:9090"
  description = "URL to access Prometheus"
}
