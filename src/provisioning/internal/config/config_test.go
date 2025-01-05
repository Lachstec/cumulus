package config

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	os.Setenv("DB_HOST", "sample.com")
	os.Setenv("DB_PORT", "1337")
	os.Setenv("DB_USER", "sample_user")
	os.Setenv("DB_PASS", "sample_pass")
	os.Setenv("AUTH0_URL", "https://auth0.com/test")
	os.Setenv("AUTH0_AUDIENCE", "https://mc-hosting.zip")
	os.Setenv("AUTH0_SECRET", "sample_secret")
	os.Setenv("OPENSTACK_IDENTITY_ENDPOINT", "https://my-keystoneserver.zip/v3")
	os.Setenv("OPENSTACK_USER", "stackymcstackface")
	os.Setenv("OPENSTACK_PASS", "secure_password1!")
	os.Setenv("OPENSTACK_DOMAIN", "default")
	os.Setenv("OPENSTACK_TENANT_NAME", "my_tenant")
	os.Setenv("CRYPTO_KEY", "my_secret_key")

	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("DB_USER")
		os.Unsetenv("DB_PASS")
		os.Unsetenv("AUTH0_URL")
		os.Unsetenv("AUTH0_AUDIENCE")
		os.Unsetenv("AUTH0_SECRET")
		os.Unsetenv("OPENSTACK_IDENTITY_ENDPOINT")
		os.Unsetenv("OPENSTACK_USER")
		os.Unsetenv("OPENSTACK_PASS")
		os.Unsetenv("OPENSTACK_DOMAIN")
		os.Unsetenv("OPENSTACK_TENANT_NAME")
		os.Unsetenv("CRYPTO_KEY")
	}()

	cfg := LoadConfig()

	assert.Equal(t, "sample.com", cfg.Db.Host)
	assert.Equal(t, "1337", cfg.Db.Port)
	assert.Equal(t, "sample_user", cfg.Db.User)
	assert.Equal(t, "sample_pass", cfg.Db.Password)
	assert.Equal(t, url.URL{Scheme: "https", Host: "auth0.com", Path: "/test"}, cfg.Auth0.AuthURL)
	assert.Equal(t, "https://my-keystoneserver.zip/v3", cfg.Openstack.IdentityEndpoint)
	assert.Equal(t, "stackymcstackface", cfg.Openstack.Username)
	assert.Equal(t, "secure_password1!", cfg.Openstack.Password)
	assert.Equal(t, "default", cfg.Openstack.Domain)
	assert.Equal(t, "my_tenant", cfg.Openstack.TenantName)
	assert.Equal(t, []byte("my_secret_key"), cfg.CryptoConfig.EncryptionKey)
	assert.Equal(t, url.URL{Scheme: "https", Host: "auth0.com", Path: "/test"}, cfg.Auth0.AuthURL)
	assert.Equal(t, "https://mc-hosting.zip", cfg.Auth0.Audience)
	assert.Equal(t, "sample_secret", cfg.Auth0.Secret)
}

func TestFallbackValues(t *testing.T) {
	cfg := LoadConfig()

	assert.Equal(t, "localhost", cfg.Db.Host)
	assert.Equal(t, "5432", cfg.Db.Port)
	assert.Equal(t, "postgres", cfg.Db.User)
	assert.Equal(t, "postgres", cfg.Db.Password)
	assert.Equal(t, url.URL{Scheme: "http", Host: "localhost"}, cfg.Auth0.AuthURL)
	assert.Equal(t, "http://localhost", cfg.Openstack.IdentityEndpoint)
	assert.Equal(t, "osuser", cfg.Openstack.Username)
	assert.Equal(t, "ospassword", cfg.Openstack.Password)
	assert.Equal(t, "osp", cfg.Openstack.Domain)
	assert.Equal(t, "default", cfg.Openstack.TenantName)
	assert.Equal(t, []byte("super_secure_default_key1!"), cfg.CryptoConfig.EncryptionKey)
	assert.Equal(t, url.URL{Scheme: "http", Host: "localhost"}, cfg.Auth0.AuthURL)
	assert.Equal(t, "http://localhost", cfg.Auth0.Audience)
	assert.Equal(t, "secret", cfg.Auth0.Secret)
}
