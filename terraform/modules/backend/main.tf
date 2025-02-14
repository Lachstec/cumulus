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
  name       = "backend_key"
  public_key = tls_private_key.generated.public_key_openssh
}

# Create a dedicated network for the backend
resource "openstack_networking_network_v2" "backend_network" {
  name = var.backend_network_name
}

resource "openstack_networking_subnet_v2" "backend_subnet" {
  name            = var.backend_subnet_name
  network_id      = openstack_networking_network_v2.backend_network.id
  cidr            = var.backend_subnet_cidr
  gateway_ip      = "10.10.10.17"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

# Create a dedicated security group for the backend
resource "openstack_networking_secgroup_v2" "backend_secgroup" {
  name = var.backend_security_group_name
}

# Ingress rule to allow traffic from the frontend subnet
resource "openstack_networking_secgroup_rule_v2" "backend_ingress_from_frontend" {
  security_group_id = openstack_networking_secgroup_v2.backend_secgroup.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = var.backend_port
  port_range_max    = var.backend_port
  remote_ip_prefix  = "0.0.0.0/0"
}

# Egress rule to allow traffic to the database pool
resource "openstack_networking_secgroup_rule_v2" "backend_egress_to_db" {
  security_group_id = openstack_networking_secgroup_v2.backend_secgroup.id
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = var.backend_db_port
  port_range_max    = var.backend_db_port
  remote_ip_prefix  = var.backend_db_cidr
}

# Egress rule to allow traffic to Loki
#resource "openstack_networking_secgroup_rule_v2" "backend_egress_to_loki" {
#  security_group_id = openstack_networking_secgroup_v2.backend_secgroup.id
#  direction         = "egress"
#  ethertype         = "IPv4"
#  protocol          = "tcp"
#  port_range_min    = var.loki_port
#  port_range_max    = var.loki_port
#  remote_ip_prefix  = var.loki_subnet_cidr
#}

resource "openstack_networking_secgroup_rule_v2" "backend_allow_outbound_internet" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.backend_secgroup.id
}

resource "openstack_networking_secgroup_rule_v2" "backend_allow_dns" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.backend_secgroup.id
}

resource "openstack_networking_secgroup_rule_v2" "backend_allow_icmp" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "icmp"
  security_group_id = openstack_networking_secgroup_v2.backend_secgroup.id
}

#resource "openstack_networking_router_v2" "backend_router" {
#  name = var.backend_router_name
#  external_network_id = "6f530989-999a-49e6-9197-8a33ae7bfce7"
#}

#resource "openstack_networking_router_interface_v2" "router_interface" {
#    router_id = openstack_networking_router_v2.backend_router.id
#    subnet_id = openstack_networking_subnet_v2.backend_subnet.id
#}

resource "openstack_networking_port_v2" "backend_ports" {
  count      = 2
  name       = "backend-port-${count.index + 1}"
  network_id = openstack_networking_network_v2.backend_network.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.backend_subnet.id
  }

  security_group_ids = [openstack_networking_secgroup_v2.backend_secgroup.id]
}

resource "openstack_objectstorage_container_v1" "loki_container" {
  name = "loki-logs"
}

resource "openstack_networking_secgroup_rule_v2" "loki_ingress" {
  security_group_id = openstack_networking_secgroup_v2.backend_secgroup.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 3100
  port_range_max    = 3100
  remote_ip_prefix  = "0.0.0.0/0"  # You may restrict this range for added security.
}

resource "openstack_networking_port_v2" "loki_ports" {
  count = 1
  name = "loki-port-${count.index + 1}"
  network_id = openstack_networking_network_v2.backend_network.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.backend_subnet.id
  }

  security_group_ids = [openstack_networking_secgroup_v2.backend_secgroup.id]
}

resource "openstack_compute_instance_v2" "loki_servers" {
  count = 1
  name = "loki-server-${count.index + 1}"
  image_id = var.backend_image_id
  flavor_id = var.backend_flavor_id
  key_pair = openstack_compute_keypair_v2.keypair.name

  network {
    port = openstack_networking_port_v2.loki_ports[count.index].id
  }

  user_data = templatefile("./${path.module}/loki-init.sh.tpl", {
    openstack_auth_url = var.openstack_auth_url
    openstack_username = var.openstack_username
    openstack_password = var.openstack_password
    openstack_tenant_name = var.openstack_tenant
    openstack_domain_name = var.openstack_domain_name
    openstack_region_name = var.openstack_region
    loki_container_name = openstack_objectstorage_container_v1.loki_container.name
  })
}

resource "openstack_networking_secgroup_rule_v2" "grafana_ingress" {
  security_group_id = openstack_networking_secgroup_v2.backend_secgroup.id
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 3000
  port_range_max    = 3000
  remote_ip_prefix  = "0.0.0.0/0"  # Restrict this as needed for your security requirements.
}


resource "openstack_compute_instance_v2" "backend_servers" {
  count     = 2
  name      = "backend-server-${count.index + 1}"
  image_id  = var.backend_image_id
  flavor_id = var.backend_flavor_id
  key_pair  = openstack_compute_keypair_v2.keypair.name

  network {
    port = openstack_networking_port_v2.backend_ports[count.index].id
  }

  user_data = templatefile("./${path.module}/backend-init.sh.tpl", {
    DB_HOST                     = var.backend_db_host
    DB_PORT                     = var.backend_db_port
    DB_USER                     = var.backend_db_user
    DB_PASS                     = var.backend_db_password
    OPENSTACK_IDENTITY_ENDPOINT = var.openstack_auth_url
    OPENSTACK_USER              = var.openstack_username
    OPENSTACK_PASS              = var.openstack_password
    OPENSTACK_DOMAIN            = var.openstack_domain_name
    OPENSTACK_TENANT_NAME       = var.openstack_tenant
    CRYPTO_KEY                  = var.backend_crypto_key
    TRACE_ENDPOINT              = format("http:%s:3100/loki/api/v1/push", openstack_compute_instance_v2.loki_servers[0].access_ip_v4)
    TRACE_SERVICENAME           = var.backend_tracing_service_name
    AUTH0_URL                   = var.backend_auth0_url
    AUTH0_SECRET                = var.backend_auth0_clientid
    AUTH0_AUDIENCE              = var.backend_auth0_audience
  })
}

resource "openstack_lb_loadbalancer_v2" "backend_loadbalancer" {
  name          = var.backend_loadbalancer_name
  vip_subnet_id = openstack_networking_subnet_v2.backend_subnet.id
  description   = "Loadbalancer for Backend services"
}

resource "openstack_lb_listener_v2" "backend_loadbalancer_listener" {
  name            = format("%s/%s", var.backend_loadbalancer_name, "_listener")
  protocol        = "TCP"
  protocol_port   = var.backend_port
  loadbalancer_id = openstack_lb_loadbalancer_v2.backend_loadbalancer.id
}

resource "openstack_lb_pool_v2" "backend_loadbalancer_pool" {
  name        = format("%s/%s", var.backend_loadbalancer_name, "_pool")
  protocol    = "TCP"
  lb_method   = "ROUND_ROBIN"
  listener_id = openstack_lb_listener_v2.backend_loadbalancer_listener.id
}

resource "openstack_lb_member_v2" "backend_loadbalancer_members" {
  count         = 2
  pool_id       = openstack_lb_pool_v2.backend_loadbalancer_pool.id
  address       = openstack_compute_instance_v2.backend_servers[count.index].access_ip_v4
  protocol_port = var.backend_port
  subnet_id     = openstack_networking_subnet_v2.backend_subnet.id
}

resource "openstack_lb_monitor_v2" "backend_loadbalancer_healthcheck" {
  pool_id     = openstack_lb_pool_v2.backend_loadbalancer_pool.id
  type        = "TCP"
  delay       = 5
  timeout     = 5
  max_retries = 3
}