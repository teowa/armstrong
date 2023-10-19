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

resource "azapi_resource" "trafficManagerProfile" {
  type      = "Microsoft.Network/trafficManagerProfiles@2018-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = jsonencode({
    properties = {
      dnsConfig = {
        relativeName = "acctest-tmp-230630034107608613"
        ttl          = 30
      }
      monitorConfig = {
        expectedStatusCodeRanges = [
        ]
        intervalInSeconds         = 30
        path                      = "/"
        port                      = 443
        protocol                  = "HTTPS"
        timeoutInSeconds          = 10
        toleratedNumberOfFailures = 3
      }
      trafficRoutingMethod = "Weighted"
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "ExternalEndpoint" {
  type      = "Microsoft.Network/trafficManagerProfiles/ExternalEndpoints@2018-08-01"
  parent_id = azapi_resource.trafficManagerProfile.id
  name      = var.resource_name
  body = jsonencode({
    properties = {
      customHeaders = [
      ]
      endpointStatus = "Enabled"
      subnets = [
      ]
      target = "www.example.com"
      weight = 3
    }
  })
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

