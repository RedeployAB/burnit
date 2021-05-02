output "fqdn" {
  value       = azurerm_container_group.container_group.fqdn
  description = "FQDN of exposed container"
}

output "ip_address" {
  value       = azurerm_container_group.container_group.ip_address
  description = "Public IP of exposed container"
}
