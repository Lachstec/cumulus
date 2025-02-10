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
  backend_cidr = "10.10.10.16/28"
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