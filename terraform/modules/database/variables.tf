##############################
# PostgreSQL Instance Variables
##############################
variable "postgres_image_id" {
  description = "Image ID for postgres database nodes"
  type        = string
}

variable "postgres_flavor_id" {
  description = "Flavor ID for postgres database nodes"
  type        = string
}

##############################
# pgpool Instance Variables
##############################
variable "pgpool_image_id" {
  description = "Image ID for pgpool node"
  type        = string
}

variable "pgpool_flavor_id" {
  description = "Flavor ID for pgpool node"
  type        = string
}

##############################
# Network Variables for Dedicated Internal Network
##############################
variable "postgres_network_name" {
  description = "Name for the dedicated PostgreSQL/pgpool network"
  type        = string
  default     = "postgres-network"
}

variable "postgres_subnet_name" {
  description = "Name for the dedicated PostgreSQL/pgpool subnet"
  type        = string
  default     = "postgres-subnet"
}

variable "postgres_subnet_cidr" {
  description = "CIDR for the dedicated PostgreSQL/pgpool subnet"
  type        = string
  default     = "10.10.10.0/28"
}

##############################
# Security Groups
##############################
variable "postgres_security_group_name" {
  description = "Name of the security group for postgres instances"
  type        = string
  default     = "postgres-sg"
}

variable "pgpool_security_group_name" {
  description = "Name of the security group for pgpool instances"
  default     = "pgpool-sg"
}

##############################
# Database User
##############################
variable "postgres_user" {
  description = "Name of the postgres user that should get created"
  type        = string
}

variable "postgres_password" {
  description = "Password of the postgres user that should get created"
  type        = string
}

variable "pgpool_user" {
  type        = string
  description = "The user for Pgpool health checks"
}

variable "pgpool_password" {
  type        = string
  description = "The password for Pgpool health checks"
}

variable "backend_cidr" {
  description = "CIDR for the backend subnet"
  type = string
}