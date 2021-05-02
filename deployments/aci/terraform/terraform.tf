terraform {
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "=2.57.0"
    }
  }
}

provider "azurerm" {
  disable_terraform_partner_id = true
  storage_use_azuread          = true
  features {}
}