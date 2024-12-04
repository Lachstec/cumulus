package config

import (
	"github.com/joho/godotenv"
	"log"
	"net/url"
	"os"
)

// Config is the main configuration type for the application.
// All settings or credentials the backend needs is to be supplied here.
type Config struct {
	// Db configuration to access the primary database
	Db    DbConfig
	Auth0 Auth0Config
	// Openstack configuration to connect to an openstack cluster.
	Openstack OpenStackConfig
}

// LoadConfig loads the application configuration.
//
// It first checks if a .env file is available and loads it if so, overriding possibly set
// environment variables. If it is not present, it simply reads the set environment
// variables, supplying a default value if one should not be present.
// Expected variables:
//
// DB_HOST: hostname for the database (default: localhost)
// DB_PORT: port for the database (default: 5432)
// DB_USER: username for the database (default: postgres)
// DB_PASS: password for the database (default: postgres)
// AUTH0_URL: URL to the Auth0 Userinfo endpoint (default: http://localhost)
// OPENSTACK_IDENTITY_ENDPOINT: Keystone URL of the openstack cluster (default: http://localhost)
// OPENSTACK_USER: Username for openstack (default: osuser)
// OPENSTACK_PASS: Password for openstack (default: ospassword)
// OPENSTACK_TENANT_ID: TenantID for openstack (default: osp)
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using Fallback values")
	}

	auth0, err := url.Parse(getEnv("AUTH0_URL", "http://localhost"))
	if err != nil {
		log.Fatalln("Invalid url for Auth0")
	}

	cfg := &Config{
		Db: DbConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASS", "postgres"),
		},
		Auth0: Auth0Config{
			Url: *auth0,
		},
		Openstack: OpenStackConfig{
			identityEndpoint: getEnv("OPENSTACK_IDENTITY_ENDPOINT", "http://localhost"),
			username:         getEnv("OPENSTACK_USER", "osuser"),
			password:         getEnv("OPENSTACK_PASS", "ospassword"),
			tenantId:         getEnv("OPENSTACK_TENANT_ID", "osp"),
		},
	}
	return cfg
}

// getEnv looks if an environment variable with the name key exists,
// returning it if so, else, the fallback is returned.
func getEnv(key string, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
