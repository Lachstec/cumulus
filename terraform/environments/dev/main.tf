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
  backend_cidr       = "10.10.10.16/28"
}

module "backend" {
  source = "../../modules/backend"

  providers = {
    openstack = openstack
    tls = tls
  }

  depends_on = [ module.database ]

  # Backend instance configuration
  backend_image_id       = "d6d1835c-7180-4ca9-b4a1-470afbd8b398"
  backend_flavor_id      = "3"
  backend_network_name   = "backend-network"
  backend_subnet_name    = "backend-subnet"
  backend_subnet_cidr    = "10.10.10.16/28"
  backend_security_group_name = "backend-sg"
  backend_router_name    = "backend-router"
  backend_loadbalancer_name = "backend-lb"

  backend_db_host        = module.database.pgpool_ip
  backend_db_port        = "9999"
  backend_db_user        = var.pgpool_user
  backend_db_password    = var.pgpool_password
  backend_db_cidr        = module.database.pg_subnet_cidr

  # OpenStack authentication
  openstack_auth_url     = "https://10.32.7.184:5000/v3"
  openstack_region       = "nova"
  openstack_username     = "CloudServ11"
  openstack_password     = "demo"
  openstack_tenant       = "CloudServ11"
  openstack_domain_name  = "Default"

  # Encryption & Security
  backend_crypto_key     = "1YRCJE3rUygZv4zXUhBNUf1sDUIszdT2KAtczVYB85c="

  # Logging & Monitoring
  backend_tracing_endpoint   = "http://tracing.example.com:4317"
  backend_tracing_service_name = "backend-service"

  # API configuration
  backend_port           = "10000"

  # Frontend configuration
  frontend_subnet_cidr   = "10.10.30.0/24"

  # Auth0 configuration
  backend_auth0_url      = "https://your-tenant.auth0.com/"
  backend_auth0_clientid = "your-client-id"
  backend_auth0_audience = "your-audience"
}

module "frontend" {
  source = "../../modules/frontend"

  providers = {
    openstack = openstack
  }

  frontend_client_id = module.auth0.frontend_client_id
  frontend_flavor_id = "3"
  frontend_image_id = "d6d1835c-7180-4ca9-b4a1-470afbd8b398"
  frontend_backend_url = format("http://%s:%d", module.floating_ips.backend_lb_floating_ip, 10000)
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