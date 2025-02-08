variable "openstack_auth_url" {
    description = "The OpenStack authentication URL"
    type = string
}

variable "openstack_region" {
    description = "The OpenStack region to use"
    type = string
}

variable "openstack_username" {
    description = "OpenStack username to use for deployment"
    type = string
}

variable "openstack_password" {
    description = "OpenStack password to use for deployment"
    type = string
    sensitive = true
}

variable "openstack_tenant" {
    description = "OpenStack tenant or project name"
    type = string
}

variable "openstack_domain_name" {
    description = "OpenStack domain name"
    type = string
    default = "Default"
}