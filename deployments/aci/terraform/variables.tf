# Resource group variables
#############################
variable "resource_group_name" {
  type = string
}

variable "resource_group_deploy" {
  type    = bool
  default = true
}

# All resources variables
#############################
variable "location" {
  type    = string
  default = "westeurope"
}

variable "tags" {
  type    = map(string)
  default = {}
}

# Redis variables
#############################
variable "redis_cache_name" {
  type = string
}

variable "redis_cache_deploy" {
  type    = bool
  default = true
}

# ACR and registry variables
#############################
variable "acr" {
  type    = bool
  default = true
}

variable "acr_name" {
  type = string
}

variable "acr_resource_group_name" {
  type = string
}

variable "registry" {
  type    = string
  default = ""
}

variable "registry_username" {
  type    = string
  default = ""
}

variable "registry_password" {
  type    = string
  default = ""
}

# Container group variables
#############################
variable "container_group_name" {
  type    = string
  default = "burnit"
}

variable "ip_address_type" {
  type    = string
  default = "Public"
}

variable "restart_policy" {
  type    = string
  default = "Always"
}

# Service variables
#############################
variable "nginx" {
  type = object({
    name    = string
    version = string
    port    = number
    cpu     = string
    memory  = string
  })
  default = {
    name    = "nginx"
    version = "1.20.0-alpine"
    port    = 443
    cpu     = "0.5"
    memory  = "0.5"
  }
}

variable "gw" {
  type = object({
    name    = string
    version = string
    port    = number
    cpu     = string
    memory  = string
  })
  default = {
    name    = "burnitgw"
    version = "0.1.0"
    port    = 3000
    cpu     = "0.5"
    memory  = "0.5"
  }
}

variable "db" {
  type = object({
    name    = string
    version = string
    port    = number
    cpu     = string
    memory  = string
  })
  default = {
    name    = "burnitdb"
    version = "0.1.0"
    port    = 3001
    cpu     = "0.5"
    memory  = "0.5"
  }
}

variable "gen" {
  type = object({
    name    = string
    version = string
    port    = number
    cpu     = string
    memory  = string
  })
  default = {
    name    = "burnitgen"
    version = "0.1.0"
    port    = 3002
    cpu     = "0.5"
    memory  = "0.5"
  }
}

variable "ssl_certificate_path" {
  type = string
}

variable "ssl_key_path" {
  type = string
}

variable "nginx_config_path" {
  type = string
}

variable "encryption_key" {
  type      = string
  sensitive = true
  default   = ""
}
