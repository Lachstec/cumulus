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

provider "openstack" {
  auth_url     = "#"
  region       = "#"
  password     = "#"
  domain_name  = "#"
  user_name     = "#"
  tenant_name   = "#"
  insecure = true
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

resource "openstack_lb_loadbalancer_v2" "frontend_loadbalancer" {
  name = var.frontend_loadbalancer_name
  vip_subnet_id = openstack_networking_subnet_v2.frontend_subnet.id
  description = "Loadbalancer for frontend services"
}

resource "openstack_lb_listener_v2" "frontend_loadbalancer_listener" {
  name = format("%s/%s",var.frontend_loadbalancer_name,"_listener")
  protocol = "TCP"
  protocol_port = var.frontend_port
  loadbalancer_id = openstack_lb_loadbalancer_v2.frontend_loadbalancer.id
}

resource "openstack_lb_pool_v2" "frontend_loadbalancer_pool" {
  name = format("%s/%s",var.frontend_loadbalancer_name,"_pool")
    protocol = "TCP"
    lb_method = "ROUND_ROBIN"
    listener_id = openstack_lb_listener_v2.frontend_loadbalancer_listener.id
}

resource "openstack_lb_member_v2" "frontend_loadbalancer_members" {
  count             = 2
  pool_id           = openstack_lb_pool_v2.frontend_loadbalancer_pool.id
  address           = openstack_compute_instance_v2.frontend_servers[count.index].access_ip_v4
  protocol_port     = var.frontend_port
  subnet_id         = openstack_networking_subnet_v2.frontend_subnet.id
}

resource "openstack_lb_monitor_v2" "frontend_loadbalancer_healthcheck" {
  pool_id = openstack_lb_pool_v2.frontend_loadbalancer_pool.id
  type = "HTTP"
  delay             = 5
  timeout           = 5
  max_retries       = 3
}

# Neeed a router for the frontend
resource "openstack_networking_router_v2" "frontend_router" {
  name = "frontend-router"
  external_network_id = "6f530989-999a-49e6-9197-8a33ae7bfce7"
}

resource "openstack_networking_router_interface_v2" "router_interface" {
    router_id = openstack_networking_router_v2.frontend_router.id
    subnet_id = openstack_networking_subnet_v2.frontend_subnet.id
}

# Creating floating IP for the frontend - to test. Should happen in Main before Auth0
resource "openstack_networking_floatingip_v2" "frontend_floatingip" {
  pool = "ext_net"
  description = "frontend floating ip"
  port_id = openstack_lb_loadbalancer_v2.frontend_loadbalancer.vip_port_id
}
# Here we could create the Auth0 service, and then create the frontend itself