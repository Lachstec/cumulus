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
  }
}

resource "openstack_networking_floatingip_v2" "backend_lb_floating_ip" {
  pool        = var.external_network_name
  description = "cumulus_backend"
}

resource "openstack_networking_floatingip_v2" "frontend_lb_floating_ip" {
  pool        = var.external_network_name
  description = "cumulus_frontend"
}