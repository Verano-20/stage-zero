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

  lifecycle {
    prevent_destroy = true
    ignore_changes = []
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