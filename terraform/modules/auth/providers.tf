terraform {
  required_version = ">= 1.5.0"
  required_providers {
    auth0 = {
      source  = "auth0/auth0"
      version = ">= 1.0.0"
    }
  }
}

provider "auth0" {
  domain        = var.auth_domain
  client_id     = var.auth_client_id
  client_secret = var.auth_client_secret
}