package config

import "net/url"

// Auth0Config represents Configuration related to connecting to an Auth0 Tenant
type Auth0Config struct {
	AuthURL  url.URL
	Audience string
	Secret   string
}
