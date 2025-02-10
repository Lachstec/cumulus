output "pgpool_ip" {
  description = "Access IP address of the pgpool instance"
  value       = openstack_compute_instance_v2.pgpool.access_ip_v4
}

output "pg_ips" {
  description = "List of PostgreSQL instance IPs"
  value       = [for instance in openstack_compute_instance_v2.pgsql : instance.access_ip_v4]
}

output "pg_network_id" {
  description = "ID of the dedicated PostgreSQL network"
  value       = openstack_networking_network_v2.pg_network.id
}

output "pg_subnet_id" {
  description = "ID of the dedicated PostgreSQL subnet"
  value       = openstack_networking_subnet_v2.pg_subnet.id
}