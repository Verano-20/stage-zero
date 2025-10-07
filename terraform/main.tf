locals {
    droplet_name = "stage-zero-tf"
    region = "lon1"
}

resource "digitalocean_droplet" "stage-zero" {
  name   = local.droplet_name
  image  = "ubuntu-25-04-x64"
  size   = "s-1vcpu-1gb"
  region = local.region
  ssh_keys = [data.digitalocean_ssh_key.terraform.id]
  tags   = ["terraform", local.droplet_name]

  user_data = templatefile("../scripts/user-data.sh", {
    droplet_name = local.droplet_name
    github_token = var.github_token
    github_username = var.github_username
  })

  lifecycle {
    prevent_destroy = false
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
