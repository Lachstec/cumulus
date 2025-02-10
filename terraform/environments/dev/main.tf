module "database" {
  source = "../../modules/database"

  postgres_image_id  = var.postgres_image_id
  postgres_flavor_id = var.postgres_flavor_id
  pgpool_image_id    = var.pgpool_image_id
  pgpool_flavor_id   = var.pgpool_flavor_id
  postgres_user      = var.postgres_user
  postgres_password  = var.postgres_password
  pgpool_user        = var.pgpool_user
  pgpool_password    = var.pgpool_password
}

module "backend" {
  source = "../../modules/backend"
}