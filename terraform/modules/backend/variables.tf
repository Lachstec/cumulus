variable "backend_image_id" {
  description = "The Cloud Image to use for the Backend Instance"
  type        = string
}

variable "backend_flavor_id" {
  description = "The Flavor ID to use for the Backend Instance"
  type        = string
}

variable "backend_network_name" {
  description = "Name of the network for the Backend servers"
  type        = string
  default     = "backend-network"
}

variable "backend_subnet_name" {
  description = "Name of the subnet dedicated for Backend services"
  type        = string
  default     = "backend-subnet"
}

variable "backend_subnet_cidr" {
  description = "CIDR Notation of backend subnet address and mask"
  type        = string
  default     = "10.10.10.16/28"
}

variable "backend_security_group_name" {
  description = "Name of the security group for the backend"
  type        = string
  default     = "backend-sg"
}

variable "backend_router_name" {
    description = "Name of the router for the backend"
    type = string
}

variable "backend_loadbalancer_name" {
  description = "Name of the loadbalancer for the backend"
  type = string
}

variable "backend_db_host" {
  description = "Hostname where a PostgreSQL Database is reachable for the backend"
  type        = string
}

variable "backend_db_port" {
  description = "Port of the PostgreSQL Database"
  type        = string
}

variable "backend_db_user" {
  description = "Username of the Database user to use"
  type        = string
}

variable "backend_db_password" {
  description = "Password of the Database user to use"
  type        = string
  sensitive   = true
}

variable "backend_db_cidr" {
  description = "CIDR Notation of Database subnet"
  type = string
}

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

variable "backend_crypto_key" {
  description = "Key used to encrypt sensitive information in the database. Must be Base64 encoded and of proper length for AES-256"
  type        = string
  sensitive   = true
}

variable "backend_tracing_endpoint" {
  description = "Endpoint where Logs can be sent to"
  type        = string
}

variable "backend_tracing_service_name" {
  description = "The service name to use when sending logs"
  type        = string
}

variable "backend_port" {
  description = "Port under which the API listens for connections"
  type = string
}

variable "frontend_subnet_cidr" {
    description = "The Subnet CIDR of the frontend subnet"
    type = string
}

variable "backend_auth0_url" {
  description = "URL for the Auth0 tenant"
  type = string
}

variable "backend_auth0_clientid" {
  description = "ClientID to use for Auth0"
  type = string
  sensitive = true
}

variable "backend_auth0_audience" {
  description = "Audience to use for Auth0"
  type = string
}