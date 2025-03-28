#!/usr/bin/env bash

# This scipt sets OpenStack RC environment variables as TF_VARs for Terraform
# Usage: source envPrep.sh

export TF_VAR_openstack_auth_url=$OS_AUTH_URL
export TF_VAR_openstack_region=$OS_REGION_NAME
export TF_VAR_openstack_username=$OS_USERNAME
export TF_VAR_openstack_password=$OS_PASSWORD
export TF_VAR_openstack_tenant=$OS_PROJECT_NAME