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

# Create a Keypair for the backend instances
resource "tls_private_key" "generated" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "openstack_compute_keypair_v2" "keypair" {
  name       = "frontend_key"
  public_key = tls_private_key.generated.public_key_openssh
}

# Create a dedicated network for the backend
resource "openstack_networking_network_v2" "frontend_network" {
  name = var.frontend_network_name
}

resource "openstack_networking_subnet_v2" "frontend_subnet" {
  name = var.frontend_subnet_name
  network_id = openstack_networking_network_v2.frontend_network.id
  cidr = var.frontend_subnet_cidr
  ip_version = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# Create a dedicated security group for the frontend
resource "openstack_networking_secgroup_v2" "frontend_secgroup" {
  name = var.frontend_security_group_name
}

resource "openstack_networking_port_v2" "frontend_ports" {
  count = 2
  name = "frontend-port-${count.index + 1}"
  network_id = openstack_networking_network_v2.frontend_network.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.frontend_subnet.id
  }

  security_group_ids = [ openstack_networking_secgroup_v2.frontend_secgroup.id ]
}

resource "openstack_compute_instance_v2" "frontend_servers" {
  count = 2
  name = "frontend-server-${count.index + 1}"
  image_id = var.frontend_image_id
  flavor_id = var.frontend_flavor_id
  key_pair = openstack_compute_keypair_v2.keypair.name

  network {
    port = openstack_networking_port_v2.frontend_ports[count.index].id
  }

  user_data = templatefile("./${path.module}/frontend-init.sh.tpl", {
    PUBLIC_AUTH_DOMAIN          = var.frontend_auth_url 
    PUBLIC_AUTH_CLIENT_ID       = var.frontend_client_id
    PUBLIC_AUTH_AUDIENCE        = var.frontend_audience
    PUBLIC_AUTH_CACHE_LOCATION  = var.frontend_cache_location
    PUBLIC_BACKEND_URL          = var.frontend_backend_url
    PUBLIC_REQUESTER_NAME       = var.frontend_requester_name
  })
}