output "backend_loadbalancer_ip" {
    description = "IP of the loadbalancer for the backend"
    value = openstack_lb_loadbalancer_v2.backend_loadbalancer.vip_address
}

output "backend_subnet_cidr" {
    description = "The subnet where the backend resides"
    value  = openstack_networking_subnet_v2.backend_subnet.cidr
}

output "backend_server_ips" {
  value = [for server in openstack_compute_instance_v2.backend_servers : server.access_ip_v4]
  description = "The IP addresses of the backend servers"
}