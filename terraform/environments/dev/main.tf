resource "openstack_networking_router_v2" "backend_post_router" {
  name = "backend-post-router"
  external_network_id = "6f530989-999a-49e6-9197-8a33ae7bfce7"
}

module "floating_ips" {
  source = "../../modules/floating_ips"

  providers = {
    openstack = openstack
  }

  external_network_name = var.external_network_name
}

module "auth0" {
  source = "../../modules/auth"

  providers = {
    auth0 = auth0
  }

  frontend_url = format("http://%s", module.floating_ips.frontend_lb_floating_ip)
  backend_url = format("http://%s", module.floating_ips.backend_lb_floating_ip)
  auth_client_id = var.auth_client_id
  auth_client_secret = var.auth_client_secret
  auth_domain = var.auth_domain
}

module "database" {
  source = "../../modules/database"

  providers = {
    openstack = openstack
    tls = tls
  }

  postgres_image_id  = var.postgres_image_id
  postgres_flavor_id = var.postgres_flavor_id
  pgpool_image_id    = var.pgpool_image_id
  pgpool_flavor_id   = var.pgpool_flavor_id
  postgres_user      = var.postgres_user
  postgres_password  = var.postgres_password
  pgpool_user        = var.pgpool_user
  pgpool_password    = var.pgpool_password
  backend_cidr       = var.backend_subnet_cidr
}

module "backend" {
  source = "../../modules/backend"

  providers = {
    openstack = openstack
    tls = tls
  }

  depends_on = [ module.database ]

  # Backend instance configuration
  backend_image_id       = var.backend_image_id
  backend_flavor_id      = var.backend_flavor_id
  backend_network_name   = var.backend_network_name
  backend_subnet_name    = var.backend_subnet_name
  backend_subnet_cidr    = var.backend_subnet_cidr
  backend_security_group_name = var.backend_security_group_name
  backend_router_name    = var.backend_router_name
  backend_loadbalancer_name = var.backend_loadbalancer_name

  backend_db_host        = module.database.pgpool_ip
  backend_db_port        = var.backend_db_port
  backend_db_user        = var.pgpool_user
  backend_db_password    = var.pgpool_password
  backend_db_cidr        = module.database.pg_subnet_cidr

  # OpenStack authentication
  openstack_auth_url     = var.openstack_auth_url
  openstack_region       = var.openstack_region
  openstack_username     = var.openstack_username
  openstack_password     = var.openstack_password
  openstack_tenant       = var.openstack_tenant
  openstack_domain_name  = var.openstack_domain_name

  # Encryption & Security
  backend_crypto_key     = var.backend_crypto_key

  # Logging & Monitoring
  backend_tracing_endpoint   = var.backend_tracing_endpoint
  backend_tracing_service_name = var.backend_tracing_service_name

  # API configuration
  backend_port           = var.backend_port

  # Frontend configuration
  frontend_subnet_cidr   = var.frontend_subnet_cidr

  # Auth0 configuration
  backend_auth0_url      = var.backend_auth0_url
  backend_auth0_clientid = var.backend_auth0_clientid
  backend_auth0_audience = var.backend_auth0_audience
}

module "frontend" {
  source = "../../modules/frontend"

  providers = {
    openstack = openstack
  }

  frontend_client_id = module.auth0.frontend_client_id
  frontend_flavor_id = var.frontend_flavor_id
  frontend_image_id = var.frontend_image_id
  frontend_backend_url = format("http://%s:%d", module.floating_ips.backend_lb_floating_ip, var.backend_port)
  frontend_audience = format("http://%s", module.floating_ips.backend_lb_floating_ip)
  frontend_auth_url = var.auth_domain
}


resource "openstack_networking_router_interface_v2" "router_interface" {
    router_id = openstack_networking_router_v2.backend_post_router.id
    subnet_id = module.backend.backend_subnet_id
}

resource "openstack_networking_router_interface_v2" "pg_router_interface" {
    router_id = openstack_networking_router_v2.backend_post_router.id
    subnet_id = module.database.pg_subnet_id
}

resource "openstack_networking_floatingip_associate_v2" "backend_lb_floating_ip" {
  floating_ip = module.floating_ips.backend_lb_floating_ip
  port_id = module.backend.backend_loadbalancer_vip_port
}

resource "openstack_networking_floatingip_associate_v2" "frontend_lb_floating_ip" {
  floating_ip = module.floating_ips.frontend_lb_floating_ip
  port_id = module.frontend.frontend_loadbalancer_vip_port_id
}

# Get floating IP for the load balancers
# Init auth 0
# Start Frontend