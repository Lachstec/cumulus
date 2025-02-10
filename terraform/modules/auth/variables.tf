variable "frontend_url" {
  description = "The Frontend URL, including port"
  type        = string
}

variable "backend_url" {
  description = "The Backend URL, including port"
  type        = string
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