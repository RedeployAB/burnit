resource "azurerm_resource_group" "resource_group" {
  count = var.resource_group_deploy == true ? 1 : 0

  name     = var.resource_group_name
  location = var.location
}
