package config

import (
	"encoding/base64"
	"log"
	"net/url"
	"os"

	"github.com/joho/godotenv"
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
	// TracingConfig configuration for Jaeger enabled Tracing
	TracingConfig TracingConfig
	// LoggingConfig configuration for Logging
	LoggingConfig LoggingConfig
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
// TRACE_ENDPOINT: Endpoint where logs can be sent. Contains URL and Port. (default: localhost:4317)
// TRACE_SERVICENAME: Name to pass as Servicename when sending logs to Jaeger. (default: mc-hosting)
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using Environment Variables or Fallbacks")
	}

	authURL, err := url.Parse(getEnv("AUTH0_URL", "http://localhost"))
	if err != nil {
		log.Fatalln("Invalid url for Auth0")
	}

	key, err := base64.StdEncoding.DecodeString(getEnv("CRYPTO_KEY", "1YRCJE3rUygZv4zXUhBNUf1sDUIszdT2KAtczVYB85c="))
	if err != nil {
		log.Fatalln("Invalid Encryption Key!")
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
			EncryptionKey: []byte(key),
		},
		TracingConfig: TracingConfig{
			Endpoint:    getEnv("TRACE_ENDPOINT", "localhost:4317"),
			ServiceName: getEnv("TRACE_SERVICENAME", "mc-hosting"),
		},
		LoggingConfig: LoggingConfig{
			Environment: getEnv("ENVIRONMENT", "dev"),
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
