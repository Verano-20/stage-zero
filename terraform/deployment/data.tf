data "digitalocean_ssh_key" "github_actions" {
    name = "github_actions"
}

# Find existing droplet by name
data "digitalocean_droplets" "existing" {
    filter {
        key    = "name"
        values = [local.DROPLET_NAME]
    }
}