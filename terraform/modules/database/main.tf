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
    random = {
      source  = "hashicorp/random"
      version = "~> 3.0"
    }
  }
}

# Create a Keypair for the database instances
resource "tls_private_key" "generated" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "random_password" "postgres_password" {
  length           = 16
  special          = true
  override_special = "!@#"
}

resource "openstack_compute_keypair_v2" "keypair" {
  name       = "postgres_key"
  public_key = tls_private_key.generated.public_key_openssh
}

# Create a dedicated network for the database
resource "openstack_networking_network_v2" "pg_network" {
  name = var.postgres_network_name
}

# Create a subnet for the databases
resource "openstack_networking_subnet_v2" "pg_subnet" {
  name            = var.postgres_subnet_name
  network_id      = openstack_networking_network_v2.pg_network.id
  cidr            = var.postgres_subnet_cidr
  gateway_ip      = "10.10.10.1"
  ip_version      = 4
  dns_nameservers = ["8.8.8.8", "8.8.4.4"]
}

resource "openstack_networking_secgroup_v2" "postgres_sg" {
  name        = var.postgres_security_group_name
  description = "Security group for PostgreSQL instances"
}

resource "openstack_networking_secgroup_rule_v2" "postgres_rule" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 5432
  port_range_max    = 5432
  remote_ip_prefix  = var.postgres_subnet_cidr
  security_group_id = openstack_networking_secgroup_v2.postgres_sg.id
}

resource "openstack_networking_secgroup_v2" "pgpool_sg" {
  name        = var.pgpool_security_group_name
  description = "Security group for pgpool instances"
}

resource "openstack_networking_secgroup_rule_v2" "pgpool_rule_9999" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9999
  port_range_max    = 9999
  remote_ip_prefix  = var.postgres_subnet_cidr
  security_group_id = openstack_networking_secgroup_v2.pgpool_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "pgpool_rule_5432" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 5432
  port_range_max    = 5432
  remote_ip_prefix  = var.postgres_subnet_cidr
  security_group_id = openstack_networking_secgroup_v2.pgpool_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "pgpool_egress_5432" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 5432
  port_range_max    = 5432
  remote_ip_prefix  = var.postgres_subnet_cidr
  security_group_id = openstack_networking_secgroup_v2.pgpool_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "allow_outbound_internet" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.pgpool_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "postgres_allow_outbound_internet" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 80
  port_range_max    = 443
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.postgres_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "postgres_allow_dns" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.postgres_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "pgpool_allow_dns" {
  direction         = "egress"
  ethertype         = "IPv4"
  protocol          = "udp"
  port_range_min    = 53
  port_range_max    = 53
  remote_ip_prefix  = "0.0.0.0/0"
  security_group_id = openstack_networking_secgroup_v2.pgpool_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "allow_external_pgpool" {
  direction         = "ingress"
  ethertype         = "IPv4"
  protocol          = "tcp"
  port_range_min    = 9999
  port_range_max    = 9999
  remote_ip_prefix  = var.backend_cidr
  security_group_id = openstack_networking_secgroup_v2.pgpool_sg.id
}

# Create internal ports to hook up the PostgreSQL instances
resource "openstack_networking_port_v2" "pgsql_ports" {
  count      = 2
  name       = "postgres-port-${count.index + 1}"
  network_id = openstack_networking_network_v2.pg_network.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.pg_subnet.id
  }

  security_group_ids = [openstack_networking_secgroup_v2.postgres_sg.id]
}

# Provision two PostgreSQL instances
resource "openstack_compute_instance_v2" "pgsql" {
  count     = 2
  name      = "postgres-${count.index + 1}"
  image_id  = var.postgres_image_id
  flavor_id = var.postgres_flavor_id
  key_pair  = openstack_compute_keypair_v2.keypair.name

  network {
    port = openstack_networking_port_v2.pgsql_ports[count.index].id
  }

  user_data = templatefile("${path.module}/postgres-init.sh.tpl", {
    pg_user              = "postgres"
    pg_password          = random_password.postgres_password.result
    pgpool_user          = "pgpool"
    pgpool_password      = random_password.postgres_password.result
    postgres_subnet_cidr = var.postgres_subnet_cidr
  })
}

# Create an internal port for the pgpool instance on the dedicated network
resource "openstack_networking_port_v2" "pgpool_port" {
  name       = "pgpool-port"
  network_id = openstack_networking_network_v2.pg_network.id

  fixed_ip {
    subnet_id = openstack_networking_subnet_v2.pg_subnet.id
  }

  security_group_ids = [openstack_networking_secgroup_v2.pgpool_sg.id]
}

# Provision one instance of pgpool2 for loadbalancing and HA
resource "openstack_compute_instance_v2" "pgpool" {
  name       = "pgpool"
  image_id   = var.pgpool_image_id
  flavor_id  = var.pgpool_flavor_id
  key_pair   = openstack_compute_keypair_v2.keypair.name
  depends_on = [openstack_compute_instance_v2.pgsql]

  network {
    port = openstack_networking_port_v2.pgpool_port.id
  }

  user_data = templatefile("${path.module}/pgpool-init.sh.tpl", {
    pgpool_user          = "pgpool"
    pgpool_password      = random_password.postgres_password.result
    postgres_subnet_cidr = var.postgres_subnet_cidr
    pgsql_nodes          = join(",", openstack_compute_instance_v2.pgsql[*].network[0].fixed_ip_v4)
  })
}
