terraform {
  required_providers {
    azurerm = {
      source = "hashicorp/azurerm"
      version = "=2.77.0"
    }
  }
  experiments = [ module_variable_optional_attrs ]
}

provider "azurerm"{
  features {}
}

data "azurerm_resource_group" "rg" {
  name = var.rg_name
}

data "azurerm_log_analytics_workspace" "lgworkspace" {
    name = var.loganalytics_workspace_name
    resource_group_name = var.rg_name
}

data "azurerm_subnet" "api_subnet" {
    name = var.api_subnet_name
    virtual_network_name = var.virtual_network_name
    resource_group_name = var.vnet_rg_name
}

locals {
  apim_name = "${var.uai}-${var.subcode}-apim"
  sku_name = "${var.sku}_${var.capacity}"
  api_insights = "${var.uai}-${var.subcode}-api-insights"
  api_logger = "${var.uai}-${var.subcode}-api-logger"
  api_gateway_name = "${var.uai}-${var.subcode}-gateway"
}

diagnostic_logs = [
    "GatewayLogs",
    "WebSocketConnectionLogs"
]

diagnostic_metrics = [
    "AllMetrics"
]

resource "azurerm_application_insights" "api_insights" {
  name = local.api_insights
  location = data.azurerm_resource_group.rg.location
  resource_group_name = data.azurerm_resource_group.rg.name
  application_type = "web"

  tags = {
    uai = "${var.uai}"
    env = "${var.env}"
    appname = "${var.appName}"
  }
}

resource "azurerm_api_management" "apim" {
depends_on = [azurerm_application_insights.api_insights]
name = local.apim_name
location = data.azurerm_resource_group.rg.location
resouresource_group_name = data.azurerm_resource_group.rg.name
publisher_name = "GE Gas Power"
publisher_email = "gaspowerenterprisearchitecture@ge.com"
sku_name = local.sku_name
virtual_network_type = "Internal"

virtual_network_configuration {
  subnet_id = data.azurerm_subnet.api_subnet.id
}

identity {
  type = "UserAssigned"
  identity_ids = [var.apim_identity]
}

tags = {
    uai = "${var.uai}"
    env = "${var.env}"
    appname = "${var.appName}"
}

}

#This is not required since self hosted api gateway are not approved
resource "azurerm_api_management_gateway" "apim-gateway" {
name = local.api_gateway_name
api_api_management_id = azurerm_api_management.apim.id 
description = "GE Gas Power API Management gateway"

location_data {
  name = data.azurerm_resource_group.rg.location
}
}

resource "azurerm_api_management_custom_domain" "doamin" {
  depends_on = [azurerm_api_management.apim]
  api_management_id = azurerm_api_management.apim.id

 proxy {
    host_name = "${var.appName}.power.ge.com"
    key_vault_id = var.gateway_cert_secret_id
    ssl_keyvault_identity_client_id = var.ssl_keyvault_identity_client_id
 }

 management {
   host_name = "${var.appName}-mgmt.power.ge.com"
    key_vault_id = var.mgmt_cert_secret_id
    ssl_keyvault_identity_client_id = var.ssl_keyvault_identity_client_id
 }

 developer_portal {
   host_name = "${var.appName}-portal.power.ge.com"
    key_vault_id = var.dev_portal_cert_secret_id
    ssl_keyvault_identity_client_id = var.ssl_keyvault_identity_client_id
 }

 scm {
   host_name = "${var.appName}-scm.power.ge.com"
    key_vault_id = var.scm_cert_secret_id
    ssl_keyvault_identity_client_id = var.ssl_keyvault_identity_client_id
 }
}

resource "azurerm_api_management_logger" "apim_logger" {
  depends_on = [azurerm_api_management.apim]
  name = local.api_logger
  api_management_name = azurerm_api_management.apim.name
  resource_group_name = data.azurerm_resource_group.rg.name

  application_insights {
    instrumentation_key = azurerm_application_insights.api_insights.instrumentation_key
  }
}

resource "azurerm_api_management_diagnostic" "apim_diagnostic" {
  depends_on = [azurerm_api_management.apim]
  identifier = "applicationinsights"
  resource_group_name = data.azurerm_resource_group.rg.name
  api_management_name = azurerm_api_management.apim.name
  api_management_logger_id = azurerm_api_management_logger.apim_logger.id

  sampling_percentage = 5.0
  always_log_errors = true
  log_client_ip = true
  verbosity = "verbose"
  http_correlation_protocol = "W3C"

 frontend_request {
   body_bytes = 32
   headers_to_log = [
    "content-type",
    "accept",
    "origin"
   ]
 }

 frontend_response {
   body_bytes = 32
   headers_to_log = [
    "content-type",
    "content-length",
    "origin"
   ]
 }

 backend_request {
   body_bytes = 32
   headers_to_log = [
    "content-type",
    "accept",
    "origin"
   ]
 }

 backend_response {
   body_bytes = 32
   headers_to_log = [
    "content-type",
    "content-length",
    "origin"
   ]
 }
}

resource "azurerm_monitor_diagnostic_setting" "diagnostic_setting" {
  depends_on = [azurerm_api_management.apim]
  name = "${var.uai}-${var.subcode}-${var.appName}-diagnostic"
  target_resource_id = azurerm_api_management.apim.id
  log_analytics_workspace_id = data.azurerm_log_analytics_workspace.id
  
  dynamic "log" {
    for_each = local.diagnostic_logs
    content {
        category = log.value
        retention_policy {
          enabled = false
        }
    }
  }

  dynamic "metric" {
    for_each = local.diagnostic_metrics
    content {
        category = metric.value
        retention_policy {
            enabled = false
        }
    }
  }
}