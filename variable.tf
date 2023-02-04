variable "subscription_id" {
  type = string
}

variable "subcode" {
  type = string
}

variable "rg_name" {
  type = string
}

variable "vnet_rg_name" {
  type = string
}

variable "vnet_name" {
  type = string
}

variable "uai" {
  type = string
}

variable "env" {
  type = string
  validation {
    condition = can(contains(["dev", "qa", "stg", "lab", "prd"], var.env))
    error_message = "The env variable must be one of [dev,qa,stg,lab,prd]"
  }
}

variable "sku" {
  type = string
  validation {
    condition = can(contains(["Developer", "Standard", "Premium"], var.sku))
    error_message = "The sku variable must be one of [Developer, Standard,Premium]"
  }
}

variable "capacity" {
  type = number
}

variable "appName" {
  type = string
}

variable "apim_identity" {
  type = string
}

variable "loganalytics_workspace_name" {
  type = string
}

variable "api_subnet_name" {
  type = string
}

variable "gateway_cert_secret_id" {
  type = string
}

variable "mgmt_cert_secret_id" {
  type = string
}

variable "dev_portal_cert_secret_id" {
  type = string
}

variable "portal_cert_secret_id" {
  type = string
}

variable "scm_cert_secret_id" {
  type = string
}

variable "ssl_keyvault_identity_client_id" {
  type = string
}