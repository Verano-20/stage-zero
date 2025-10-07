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
