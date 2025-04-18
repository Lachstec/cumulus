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

variable "auth_domain" {
  description = "The Auth0 domain"
  type        = string
}

variable "auth_client_id" {
  description = "The Auth0 client ID"
  type        = string
}

variable "auth_client_secret" {
  description = "The Auth0 client secret"
  type        = string
}


variable "backend_subnet_cidr" {
  description = "CIDR Notation of backend subnet address and mask"
  type        = string
  default     = "10.10.10.16/28"
}

variable "backend_db_port" {
  description = "Port of the PostgreSQL Database"
  type        = string
}

variable "backend_db_cidr" {
  description = "CIDR Notation of Database subnet"
  type        = string
}

variable "backend_crypto_key" {
  description = "Key used to encrypt sensitive information in the database. Must be Base64 encoded and of proper length for AES-256"
  type        = string
  sensitive   = true
}

variable "backend_port" {
  description = "Port under which the API listens for connections"
  type        = string
}

variable "backend_auth0_url" {
  description = "URL for the Auth0 tenant"
  type        = string
}

variable "backend_auth0_clientid" {
  description = "ClientID to use for Auth0"
  type        = string
  sensitive   = true
}

variable "backend_auth0_audience" {
  description = "Audience to use for Auth0"
  type        = string
}

variable "frontend_image_id" {
  description = "The Cloud Image to use for the Frontend Instance"
  type        = string
}

variable "frontend_flavor_id" {
  description = "The Flavor ID to use for the Frontend Instance"
  type        = string
}

variable "frontend_auth_url" {
  description = "The URL of the Auth0 service"
  type        = string
}

variable "frontend_client_id" {
  description = "The client ID for the Auth0 service"
  type        = string
}

variable "frontend_audience" {
  description = "The audience for the Auth0 service (= the Backend API)"
  type        = string
}

variable "frontend_backend_url" {
  description = "The URL of the Backend"
  type        = string
}

## Variables with sensitive values

## Variables with default values
variable "frontend_cache_location" {
  description = "The location of the cache server"
  type        = string
  default     = "localstorage"
}

variable "frontend_requester_name" {
  description = "The name of the requester"
  type        = string
  default     = "terraform-frontend"
}

variable "frontend_network_name" {
  description = "Name of the network for the Frontend servers"
  type        = string
  default     = "frontend-network"
}

variable "frontend_security_group_name" {
  description = "Name of the security group for the frontend"
  type        = string
  default     = "frontend-sg"
}

variable "frontend_subnet_name" {
  description = "Name of the subnet dedicated for Frontend services"
  type        = string
  default     = "frontend-subnet"
}

variable "frontend_subnet_cidr" {
  description = "CIDR Notation of frontend subnet address and mask"
  type        = string
  default     = "10.10.30.0/24"
}

variable "frontend_loadbalancer_name" {
  description = "Name of the loadbalancer for the frontend"
  type        = string
  default     = "frontend-lb"
}

variable "frontend_port" {
  description = "Port of the Frontend service"
  type        = string
  default     = "80"
}