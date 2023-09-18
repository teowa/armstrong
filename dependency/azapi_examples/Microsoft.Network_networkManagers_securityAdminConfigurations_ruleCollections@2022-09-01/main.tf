terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

data "azurerm_client_config" "current" {
}

data "azapi_resource" "subscription" {
  type                   = "Microsoft.Resources/subscriptions@2021-01-01"
  resource_id            = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  response_export_values = ["*"]
}

resource "azapi_resource" "resourceGroup" {
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      description = ""
      networkManagerScopeAccesses = [
        "SecurityAdmin",
      ]
      networkManagerScopes = {
        managementGroups = [
        ]
        subscriptions = [
          data.azapi_resource.subscription.id,
        ]
      }
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "securityAdminConfiguration" {
  type      = "Microsoft.Network/networkManagers/securityAdminConfigurations@2022-09-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      applyOnNetworkIntentPolicyBasedServices = []
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "networkGroup" {
  type      = "Microsoft.Network/networkManagers/networkGroups@2022-09-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "ruleCollection" {
  type      = "Microsoft.Network/networkManagers/securityAdminConfigurations/ruleCollections@2022-09-01"
  parent_id = azapi_resource.securityAdminConfiguration.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      appliesToGroups = [
        {
          networkGroupId = azapi_resource.networkGroup.id
        },
      ]
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

