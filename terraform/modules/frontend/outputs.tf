output "frontend_loadbalancer_vip_port_id" {
  description = "vip port id of the frontend loadbalancer"
  value       = openstack_lb_loadbalancer_v2.frontend_loadbalancer.vip_port_id
}