data "azurerm_container_registry" "acr" {
  count = var.acr == true ? 1 : 0

  name                = var.acr_name
  resource_group_name = var.acr_resource_group_name
}
