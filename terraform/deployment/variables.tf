variable "do_token" {
    description = "DigitalOcean API token"
    type = string
    sensitive = true
}

variable "github_token" {
    description = "GitHub token for container registry access"
    type = string
    sensitive = true
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