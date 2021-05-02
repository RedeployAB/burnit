locals {
  registry          = var.acr == true ? data.azurerm_container_registry.acr[0].login_server : var.registry
  registry_username = var.acr == true ? data.azurerm_container_registry.acr[0].admin_username : var.registry_username
  registry_password = var.acr == true ? data.azurerm_container_registry.acr[0].admin_password : var.registry_password

  nginx_image = "${local.registry}/${var.nginx.name}:${var.nginx.version}"
  gw_image    = "${local.registry}/${var.gw.name}:${var.gw.version}"
  db_image    = "${local.registry}/${var.db.name}:${var.db.version}"
  gen_image   = "${local.registry}/${var.gen.name}:${var.gen.version}"
}

resource "random_password" "encryption_key" {
  length  = 32
  special = false
}

resource "azurerm_container_group" "container_group" {
  name                = var.container_group_name
  resource_group_name = var.resource_group_name
  location            = var.location

  dns_name_label  = var.container_group_name
  ip_address_type = var.ip_address_type
  os_type         = "Linux"

  restart_policy = var.restart_policy

  image_registry_credential {
    username = local.registry_username
    password = local.registry_password
    server   = local.registry
  }

  exposed_port {
    port     = var.nginx.port
    protocol = "TCP"
  }

  container {
    name   = var.nginx.name
    image  = local.nginx_image
    cpu    = var.nginx.cpu
    memory = var.nginx.memory

    ports {
      port     = var.nginx.port
      protocol = "TCP"
    }

    volume {
      name       = "nginx-config"
      mount_path = "/etc/nginx"
      secret = {
        "ssl.crt"    = filebase64(var.ssl_certificate_path)
        "ssl.key"    = filebase64(var.ssl_key_path)
        "nginx.conf" = filebase64(var.nginx_config)
      }
    }
  }

  container {
    name   = var.gw.name
    image  = local.gw_image
    cpu    = var.gw.cpu
    memory = var.gw.memory

    ports {
      port     = var.gw.port
      protocol = "TCP"
    }
  }

  container {
    name   = var.db.name
    image  = local.db_image
    cpu    = var.db.cpu
    memory = var.db.memory

    ports {
      port     = var.db.port
      protocol = "TCP"
    }

    secure_environment_variables = {
      DB_CONNECTION_URI       = var.redis_deploy == true ? azurerm_redis_cache.redis[0].primary_connection_string : data.azurerm_redis_cache.redis[0].primary_connection_string
      BURNITDB_ENCRYPTION_KEY = var.encryption_key == "" ? random_password.encryption_key.result : var.encryption_key
    }
  }

  container {
    name   = var.gen.name
    image  = local.gen_image
    cpu    = var.gen.cpu
    memory = var.gen.memory

    ports {
      port     = var.gen.port
      protocol = "TCP"
    }
  }
}
