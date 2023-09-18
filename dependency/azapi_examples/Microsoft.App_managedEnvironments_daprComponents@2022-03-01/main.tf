terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
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

resource "azapi_resource" "resourceGroup" {
  type                      = "Microsoft.Resources/resourceGroups@2020-06-01"
  name                      = var.resource_name
  location                  = var.location
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "sharedKeys" {
  type                   = "Microsoft.OperationalInsights/workspaces@2020-08-01"
  resource_id            = azapi_resource.workspace.id
  action                 = "sharedKeys"
  response_export_values = ["*"]
}

resource "azapi_resource" "managedEnvironment" {
  type      = "Microsoft.App/managedEnvironments@2022-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = jsonencode({
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = jsondecode(azapi_resource.workspace.output).properties.customerId
          sharedKey  = jsondecode(data.azapi_resource_action.sharedKeys.output).primarySharedKey
        }
      }
      vnetConfiguration = {
      }
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "daprComponent" {
  type      = "Microsoft.App/managedEnvironments/daprComponents@2022-03-01"
  parent_id = azapi_resource.managedEnvironment.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      componentType = "state.azure.blobstorage"
      ignoreErrors  = false
      initTimeout   = "5s"
      scopes        = null
      version       = "v1"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

