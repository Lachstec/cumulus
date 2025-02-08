# Create a dedicated network for the database
resource "openstack_networking_network_v2" "pg_network" {
  name = var.postgres_network_name
}

# Create a subnet for the databases
resource "openstack_networking_subnet_v2" "pg_subnet" {
  name       = var.postgres_subnet_name
  network_id = openstack_networking_network_v2.pg_network.id
  cidr       = var.postgres_subnet_cidr
  ip_version = 4
}

resource "openstack_networking_secgroup_v2" "postgres_sg" {
    name = var.postgres_security_group_name
    description = "Security group for PostgreSQL instances"
}

resource "openstack_networking_secgroup_rule_v2" "postgres_rule" {
    direction = "ingress"
    ethertype = "IPv4"
    protocol = "tcp"
    port_range_min = 5432
    port_range_max = 5432
    remote_ip_prefix = var.postgres_subnet_cidr
    security_group_id = openstack_networking_secgroup_v2.postgres_sg.id
}

resource "openstack_networking_secgroup_v2" "pgpool_sg" {
    name = var.pgpool_security_group_name
    description = "Security group for pgpool instances"
}

resource "openstack_networking_secgroup_rule_v2" "pgpool_rule_9999" {
    direction = "ingress"
    ethertype = "IPv4"
    protocol = "tcp"
    port_range_min = 9999
    port_range_max = 9999
    remote_ip_prefix = var.postgres_subnet_cidr
    security_group_id = openstack_networking_secgroup_v2.pgpool_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "pgpool_rule_5432" {
  direction        = "ingress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 5432
  port_range_max   = 5432
  remote_ip_prefix = var.postgres_subnet_cidr
  security_group_id = openstack_networking_secgroup_v2.pgpool_sg.id
}

resource "openstack_networking_secgroup_rule_v2" "pgpool_egress_5432" {
  direction        = "egress"
  ethertype        = "IPv4"
  protocol         = "tcp"
  port_range_min   = 5432
  port_range_max   = 5432
  remote_ip_prefix = var.postgres_subnet_cidr
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
  key_pair  = var.key_name

  network {
    port = openstack_networking_port_v2.pgsql_ports[count.index].id
  }

  user_data = templatefile("${path.module}/postgres-init.sh.tpl", {
    pg_user = var.postgres_user
    pg_password = var.postgres_password
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
  name      = "pgpool"
  image_id  = var.pgpool_image_id
  flavor_id = var.pgpool_flavor_id
  key_pair  = var.key_name
  depends_on = [ openstack_compute_instance_v2.pgsql ]

  network {
    port = openstack_networking_port_v2.pgpool_port.id
  }

  provisioner "file" {
    source      = "${path.module}/failover.sh"
    destination = "/etc/pgpool2/failover.sh" 
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /etc/pgpool2/failover.sh"
    ]
  }

  provisioner "remote-exec" {
    inline = [
      "sed -i 's/^#failover_command = \".*\"/failover_command = \"/etc/pgpool2/failover.sh %d %H %P %R %r %M\"/' /etc/pgpool2/pgpool.conf",
      "sed -i 's/^#health_check_period = 0/health_check_period = 5/' /etc/pgpool2/pgpool.conf",
      "sed -i 's/^#health_check_timeout = 20/health_check_timeout = 10/' /etc/pgpool2/pgpool.conf",
      "sed -i 's/^#health_check_user = \"\"/health_check_user = \"${PGPOOL_USER}\"/' /etc/pgpool2/pgpool.conf",
      "sed -i 's/^#health_check_password = \"\"/health_check_password = \"${PGPOOL_PASSWORD}\"/' /etc/pgpool2/pgpool.conf",
      "systemctl restart pgpool2"
    ]
  }

  user_data = file("${path.module}/pgpool-init.sh")
}
