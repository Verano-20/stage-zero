variable "do_token" {
    description = "DigitalOcean API token"
    type = string
    sensitive = true
}

variable "pvt_key" {
    description = "Private key for SSH access"
    type = string
    sensitive = true
}

variable "github_token" {
    description = "GitHub token for container registry access"
    type = string
    sensitive = true
}

variable "github_username" {
    description = "GitHub username for container registry"
    type = string
}