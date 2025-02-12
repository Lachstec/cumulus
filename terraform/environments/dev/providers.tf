terraform {
  required_providers {
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.0"
    }
    openstack = {
      source  = "terraform-provider-openstack/openstack"
      version = "~> 3.0.0"
    }
    auth0 = {
      source  = "auth0/auth0"
      version = ">= 1.0.0"
    }
  }
}

provider "openstack" {
  insecure = true
}

provider "auth0" {
  domain        = var.auth_domain
  client_id     = var.auth_client_id
  client_secret = var.auth_client_secret
}