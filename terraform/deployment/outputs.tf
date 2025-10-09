output "droplet_ip" {
  value = local.current_droplet.ipv4_address
  description = "The public IP address of the droplet"
}

output "droplet_id" {
  value = local.current_droplet.id
  description = "The ID of the droplet"
}

output "droplet_name" {
  value = local.current_droplet.name
  description = "The name of the droplet"
}

output "droplet_region" {
  value = local.current_droplet.region
  description = "The region of the droplet"
}

output "droplet_size" {
  value = local.current_droplet.size
  description = "The size of the droplet"
}

output "application_url" {
  value = "http://${local.current_droplet.ipv4_address}:8080"
  description = "URL to access the Stage Zero application"
}

output "grafana_url" {
  value = "http://${local.current_droplet.ipv4_address}:3000"
  description = "URL to access Grafana dashboard (admin/admin)"
}

output "prometheus_url" {
  value = "http://${local.current_droplet.ipv4_address}:9090"
  description = "URL to access Prometheus"
}

output "debug_info" {
  value = {
    droplet_name = local.DROPLET_NAME
    existing_droplets_count = length(data.digitalocean_droplets.existing.droplets)
    current_droplet_exists = local.current_droplet != null
    current_droplet_id = local.current_droplet != null ? local.current_droplet.id : "none"
  }
  description = "Debug information about droplet detection"
}
