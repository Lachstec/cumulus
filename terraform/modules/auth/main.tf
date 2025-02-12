terraform {
  required_version = ">= 1.5.0"
  required_providers {
    auth0 = {
      source  = "auth0/auth0"
      version = ">= 1.0.0"
    }
  }
}

resource "auth0_client" "frontend" {
  name                = "Cumulus Frontend"
  description         = "Cumulus Frontend from terraform"
  app_type            = "spa"
  callbacks           = [var.frontend_url]
  allowed_logout_urls = [var.frontend_url]
  allowed_origins     = [var.frontend_url]
  web_origins         = [var.frontend_url]
  cross_origin_auth   = true

}


resource "auth0_resource_server" "backend_api" {
  name             = "Cumulus Backend API"
  identifier       = var.backend_url # Unique identifier for your API
  signing_alg      = "RS256"         # Use RS256 for better security (recommended)
  enforce_policies = true
}
