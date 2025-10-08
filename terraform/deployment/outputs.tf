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

output "droplet_region" {
  value = digitalocean_droplet.stage-zero.region
  description = "The region of the droplet"
}

output "droplet_size" {
  value = digitalocean_droplet.stage-zero.size
  description = "The size of the droplet"
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
