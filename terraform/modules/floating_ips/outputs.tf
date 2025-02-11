output "backend_lb_floating_ip" {
  description = "Floating IP of the backend loadbalancer"
  value = openstack_networking_floatingip_v2.backend_lb_floating_ip.address
}

output "backend_lb_floating_ip_id" {
  description = "ID of the Floating IP of the backend loadbalancer"
  value = openstack_networking_floatingip_v2.backend_lb_floating_ip.id
}

output "frontend_lb_floating_ip" {
  description = "Floating IP of the frontend loadbalancer"
  value = openstack_networking_floatingip_v2.frontend_lb_floating_ip.address
}

output "frontend_lb_floating_ip_id" {
  description = "ID of the Floating IP of the frontend loadbalancer"
  value = openstack_networking_floatingip_v2.frontend_lb_floating_ip.id
}