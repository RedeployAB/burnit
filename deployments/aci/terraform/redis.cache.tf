resource "azurerm_redis_cache" "redis" {
  count = var.redis_cache_deploy == true ? 1 : 0

  name                = var.redis_cache_name
  resource_group_name = var.resource_group_name
  location            = var.location

  capacity = 0
  family   = "C"
  sku_name = "Basic"

  enable_non_ssl_port = false
  minimum_tls_version = "1.2"

  redis_configuration {
  }
}

data "azurerm_redis_cache" "redis" {
  count = var.redis_cache_deploy == false ? 1 : 0

  name                = var.redis_cache_name
  resource_group_name = var.resource_group_name
}
