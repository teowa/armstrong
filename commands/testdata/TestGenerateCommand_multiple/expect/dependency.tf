
resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "server" {
  type      = "Microsoft.Sql/servers@2015-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      administratorLogin         = "acctestadmin"
      administratorLoginPassword = "t2RX8A76GrnE4EKC"
      version                    = "12.0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Sql/servers/databases@2014-04-01"
  parent_id = azapi_resource.server.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      collation                     = "SQL_LATIN1_GENERAL_CP1_CI_AS"
      createMode                    = "Default"
      maxSizeBytes                  = "268435456000"
      readScale                     = "Disabled"
      requestedServiceObjectiveName = "S0"
      zoneRedundant                 = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
