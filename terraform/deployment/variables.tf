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

variable "jwt_secret" {
    description = "JWT secret for authentication"
    type = string
    sensitive = true
}

variable "postgres_password" {
    description = "PostgreSQL container password"
    type = string
    sensitive = true
}

variable "db_password" {
    description = "Database password"
    type = string
    sensitive = true
}