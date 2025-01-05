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
	Db DbConfig
	// Auth0 configuration for authentication and authorization with Auth0
	Auth0 Auth0Config
	// Openstack configuration to connect to an openstack cluster.
	Openstack OpenStackConfig
	// CryptoConfig configuration for cryptographic components
	CryptoConfig CryptoConfig
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
// DB_USER: Username for the database (default: postgres)
// DB_PASS: Password for the database (default: postgres)
// AUTH0_URL: URL to the Auth0 Userinfo endpoint (default: http://localhost)
// OPENSTACK_IDENTITY_ENDPOINT: Keystone URL of the openstack cluster (default: http://localhost)
// OPENSTACK_USER: Username for openstack (default: osuser)
// OPENSTACK_PASS: Password for openstack (default: ospassword)
// OPENSTACK_DOMAIN Domain for openstack (default: osp)
// OPENSTACK_TENANT_NAME Tenant to use for openstack (default: default)
// CRYPTO_KEY: Key to use for encrypting SSH Keys for the game servers (default: super_secure_default_key1!)
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using Fallback values")
	}

	authURL, err := url.Parse(getEnv("AUTH0_URL", "http://localhost"))
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
			AuthURL:  *authURL,
			Audience: getEnv("AUTH0_AUDIENCE", "http://localhost"),
			Secret:   getEnv("AUTH0_SECRET", "secret"),
		},
		Openstack: OpenStackConfig{
			IdentityEndpoint: getEnv("OPENSTACK_IDENTITY_ENDPOINT", "http://localhost"),
			Username:         getEnv("OPENSTACK_USER", "osuser"),
			Password:         getEnv("OPENSTACK_PASS", "ospassword"),
			Domain:           getEnv("OPENSTACK_DOMAIN", "osp"),
			TenantName:       getEnv("OPENSTACK_TENANT_NAME", "default"),
		},
		CryptoConfig: CryptoConfig{
			EncryptionKey: []byte(getEnv("CRYPTO_KEY", "super_secure_default_key1!")),
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
