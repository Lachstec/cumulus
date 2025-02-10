variable "openstack_auth_url" {
  description = "The OpenStack authentication URL"
  type        = string
}

variable "openstack_region" {
  description = "The OpenStack region to use"
  type        = string
}

variable "openstack_username" {
  description = "OpenStack username to use for deployment"
  type        = string
}

variable "openstack_password" {
  description = "OpenStack password to use for deployment"
  type        = string
  sensitive   = true
}

variable "openstack_tenant" {
  description = "OpenStack tenant or project name"
  type        = string
}

variable "openstack_domain_name" {
  description = "OpenStack domain name"
  type        = string
  default     = "Default"
}

variable "postgres_image_id" {
  description = "Image ID for postgres database nodes"
  type        = string
}

variable "postgres_flavor_id" {
  description = "Flavor ID for postgres database nodes"
  type        = string
}

variable "pgpool_image_id" {
  description = "Image ID for pgpool node"
  type        = string
}

variable "pgpool_flavor_id" {
  description = "Flavor ID for pgpool node"
  type        = string
}

variable "postgres_user" {
  description = "Name of the postgres user that should get created"
  type        = string
}

variable "postgres_password" {
  description = "Password of the postgres user that should get created"
  type        = string
}

variable "pgpool_user" {
  description = "The user for Pgpool health checks"
  type        = string
}

variable "pgpool_password" {
  description = "The password for Pgpool health checks"
  type        = string
}