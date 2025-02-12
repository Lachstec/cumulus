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