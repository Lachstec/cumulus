##############################
# PostgreSQL Instance Variables
##############################
variable "postgres_image_id" {
  description = "Image ID for postgres database nodes"
  type        = string
  default = "d6d1835c-7180-4ca9-b4a1-470afbd8b398"
}

variable "postgres_flavor_id" {
  description = "Flavor ID for postgres database nodes"
  type        = string
  default = "3"
}

##############################
# pgpool Instance Variables
##############################
variable "pgpool_image_id" {
  description = "Image ID for pgpool node"
  type        = string
  default = "d6d1835c-7180-4ca9-b4a1-470afbd8b398"
}

variable "pgpool_flavor_id" {
  description = "Flavor ID for pgpool node"
  type        = string
  default = 3
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
  default = "pguser"
}

variable "postgres_password" {
  description = "Password of the postgres user that should get created"
  type        = string
  default = "pgpassword"
}

variable "pgpool_user" {
  type        = string
  description = "The user for Pgpool health checks"
  default = "pgpooluser"
}

variable "pgpool_password" {
  type        = string
  description = "The password for Pgpool health checks"
  default = "pgpoolpassword"
}

variable "backend_cidr" {
  description = "CIDR for the backend subnet"
  type        = string
}