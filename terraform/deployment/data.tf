data "digitalocean_ssh_key" "terraform" {
    name = "terraform"
}

# Find existing droplet by name
data "digitalocean_droplets" "existing" {
    filter {
        key    = "name"
        values = [local.DROPLET_NAME]
    }
}